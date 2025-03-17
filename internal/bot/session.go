package bot

import (
	"encoding/json"
	"errors"
	"net/http"
	"presto/internal/config"
	"presto/internal/database"
	"presto/internal/discord"
	"presto/internal/log"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type WebsocketEventPayload struct {
	Opcode         int    `json:"op"`
	Data           any    `json:"d"`
	SequenceNumber int    `json:"s"`
	Name           string `json:"t"`
}

// A session that includes the current connection with Discord's gateway, data about latency and heartbeats,
// as well as the registered commands. It also handles the events received from the websocket.
type Session struct {
	Dialer *websocket.Dialer

	// The connection to the websocket
	Connection *websocket.Conn

	// The gateway URL that will be used to connect and reconnect to the websocket
	GatewayURL string

	// The length of time in milliseconds that should be used to determine how often
	// the bot needs to send a Heartbeat event in order to maintain the active connection
	HeartbeatInterval int

	// Discord required every heartbeat to be sent along with a sequence number, which
	// is provided in most events Discord sends to the bot.
	LastSequenceNumber int

	// When the last heartbeat was sent. This will be used to measure the latency between Discord's
	// gateway and the bot
	LastHeartbeat time.Time

	// The latency between Discord's gateway and the bot
	Latency time.Duration

	Cache Cache

	RegisteredCommands []ApplicationCommandWithHandler
}

// Opens a connection with Discord's gateway following the documentation
// (https://discord.com/developers/docs/events/gateway#connection-lifecycle)
func (session *Session) Open() error {
	log.Info("Started to open a connection with Discord's gateway")
	var err error

	log.Info("Started fetching gateway data")
	gatewayData := discord.GetGateway()
	session.GatewayURL = gatewayData.URL
	log.Info("Fetched gateway data successfully")

	log.Info("Started establishing a connection with Discord's websocket")
	session.Dialer = &websocket.Dialer{}
	session.Connection, _, err = session.Dialer.Dial(session.GatewayURL, http.Header{})
	if err != nil {
		log.Fatalf("There was an error while connecting to the gateway %s: %s", session.GatewayURL, err)

		// Clear cached gateway URL to ensure a new one when reconnecting
		session.GatewayURL = ""

		// Clear cached websocket Connection
		session.Connection = nil

		return err
	}
	log.Info("Established a connection with Discord's websocket successfully")

	session.Cache = Cache{}

	log.Info("Started listening for \"Hello\" event")
	_, response, err := session.Connection.ReadMessage()
	if err != nil {
		log.Fatal("Received an error while reading/receiving \"Hello\" event: " + err.Error())
		return err
	}

	var helloEventPayload WebsocketEventPayload
	json.Unmarshal(response, &helloEventPayload)

	if helloEventPayload.Opcode != discord.HELLO_EVENT_OPCODE {
		log.Errorf("The event received after establishing a connection to Discord's gateway was not \"Hello\" (opcode 10), which goes against the documentation. Instead, received opcode %d", helloEventPayload.Opcode)
		return errors.New("event received is not \"Hello\"")
	}
	log.Info("\"Hello\" event was received successfully")

	rawHelloEventData, err := json.Marshal(helloEventPayload.Data)
	if err != nil {
		log.Error("There was an error when marshalizing the \"Hello\" event data")
		return errors.New("error when marshalizing \"Hello\" event")
	}

	var helloEventData discord.HelloEventData
	err = json.Unmarshal(rawHelloEventData, &helloEventData)
	if err != nil {
		log.Error("There was an error when unmarshalizing the \"Hello\" event data")
		return errors.New("error when unmarshalizing \"Hello\" event")
	}

	session.HeartbeatInterval = helloEventData.HeartbeatInterval

	go session.Heartbeat()
	go session.Listen(session.Connection)

	log.Info("Started sending \"Identify\" event")
	err = session.Connection.WriteJSON(WebsocketEventPayload{
		Opcode: discord.IDENTIFY_EVENT_OPCODE,
		Data: discord.IdentifyEventData{
			Token:   config.DISCORD_BOT_TOKEN,
			Intents: 1, // Currently using GUILDS intent
		},
	})
	if err == nil {
		log.Info("\"Identify\" event sent successfully")
	}

	return nil
}

func (session *Session) Heartbeat() {
	select {
	case <-time.After(time.Duration(session.HeartbeatInterval)*time.Millisecond - time.Second):
		session.SendIndividualHeartbeat()
	}
}

func (session *Session) SendIndividualHeartbeat() {
	log.Info("Started heartbeat")
	if session.LastSequenceNumber == 0 {
		session.Connection.WriteJSON(WebsocketEventPayload{
			Opcode: discord.HEARTBEAT_EVENT_OPCODE,
			Data:   session.LastSequenceNumber,
		})
		log.Info("As the application has not received any events, the data sent in the hearbeat is null")
		return
	}

	session.Connection.WriteJSON(WebsocketEventPayload{
		Opcode: discord.HEARTBEAT_EVENT_OPCODE,
	})
	log.Info("Sent heartbeat with the last sequence number %d", session.LastSequenceNumber)

	session.LastHeartbeat = time.Now().UTC()
}

func (session *Session) Listen(connection *websocket.Conn) {
	log.Info("Started listening to Discord's websocket")

	for {
		_, rawEvent, err := connection.ReadMessage()

		// If any connection was closed, this will be false
		if err != nil && session.Connection == connection {
			log.Error("There was an error while reading an event: " + err.Error())
			log.Info("The connection will be closed an there will be a reconnection attempt")

			session.Close(websocket.CloseNormalClosure)
			session.Reconnect()
			continue
		}

		var event WebsocketEventPayload
		json.Unmarshal(rawEvent, &event)

		log.Info("Received an event: %v", event)
		if event.Opcode == 0 {
			session.LastSequenceNumber = event.SequenceNumber
		}

		switch event.Opcode {
		case discord.HEARTBEAT_ACK_EVENT_OPCODE:
			session.Latency = time.Since(session.LastHeartbeat)
			log.Info("A heartbeat ACK was received and the latency is %dms", session.Latency.Milliseconds())
		case discord.RECONNECT_EVENT_OPCODE:
			session.Close(websocket.CloseNormalClosure)
			session.Reconnect()
		case discord.HEARTBEAT_EVENT_OPCODE:
			log.Info("A heartbeat event of opcode 1 was received from Discord, therefore a heartbeat will be sent immediately")
			session.SendIndividualHeartbeat()
		case discord.DISPATCH_EVENT_OPCODE:
			rawEventData, _ := json.Marshal(event.Data)
			switch event.Name {
			case discord.READY:
				log.Info("Presto is ready to go")
			case discord.INTERACTION_CREATE:
				var interactionData discord.InteractionCreatePayload
				json.Unmarshal(rawEventData, &interactionData)

				session.HandleInteractionCreateEvent(interactionData)
			case discord.GUILD_CREATE:
				// Even though an unavailable guild object might be sent,
				// it would still have the id property, which is the only one
				// that will be used
				var guildData discord.Guild
				json.Unmarshal(rawEventData, &guildData)

				result := database.Connection.Create(&database.Guild{
					ID: guildData.ID,
				})
				if result.Error != nil && !errors.Is(result.Error, gorm.ErrDuplicatedKey) {
					log.Errorf("Failed to create guild %s: %s", guildData.ID, result.Error)
					continue
				} else if result.Error != nil {
					log.Errorf("Guild %s was not created because it already exists in the database", guildData.ID)
					continue
				}

				log.Info("Created guild (%s) in database successfully\n", guildData.ID)
			}
		}
	}
}

func (session *Session) Close(code int) (err error) {
	log.Info("Started closing the websocket connection")
	err = session.Connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, ""))
	if err != nil {
		log.Error("There was an error when sending a close message")
	}

	err = session.Connection.Close()

	return err
}

func (session *Session) Reconnect() error {
	log.Info("Started attempting to reconnect")

	return session.Open()
}

func NewSession(commands []ApplicationCommandWithHandler) *Session {
	PushCommands(commands)

	return &Session{
		RegisteredCommands: commands,
	}
}
