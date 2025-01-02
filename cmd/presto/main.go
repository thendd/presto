package main

import (
	"encoding/json"
	"log"
	"time"

	"presto/internal/bot/commands/setup"
	"presto/internal/config"
	"presto/internal/constants"
	"presto/internal/handlers/api/gateway"
	ws "presto/internal/handlers/websocket"
	"presto/internal/handlers/websocket/events/interaction_create"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Started loading environment variables")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	config.LoadDiscordConfig()

	log.Println("Finished loading environment variables successfully")

	log.Println("Started fetching Discord's websocket URL")
	gatewayData := gateway.GetGateway()
	log.Println("Finished fetching Discord's WSS URL successfully")

	log.Println("Started connection to Discord's WSS")
	ws.ConnnectToWebsocket(gatewayData.URL)
	log.Println("Performed websocket handshake with Discord's WSS")

	log.Println("Started listening to WSS' messages")
	// This is the workaround I found in order to satisfy the data
	// that has to be sent to Discord on every heartbeat.
	// Using `websocket.Conn.Ping()` method works "perfectly", but I
	// don't know if it's the best approach
	var lastEvent ws.DiscordGatewayEventPayload
	incomingEvents := make(chan ws.DiscordGatewayEventPayload)

	log.Println("Triggered initial handshake")
	ws.SendIdentify()
	log.Println("Handshake ocurred successfully")

	log.Println("Started registering application commands")
	setup.RegisterCommands()
	log.Println("Finished registering commands successfully")

	go ws.ReadIncomingMessages(incomingEvents)

	for {
		select {
		case <-time.After(40 * time.Second):
			log.Println("Started sending heartbeat")
			ws.SendHeartbeat(lastEvent.SequenceNumber)
			log.Println("Sent heartbeat")
		case event := <-incomingEvents:
			lastEvent = event

			eventData, _ := json.Marshal(event.Data)

			switch event.Name {
			case constants.READY:
				log.Println("Presto is ready to go")
			case constants.INTERACTION_CREATE:
				interaction_create.HandleInteractionCreate(eventData)
			}
		}
	}
}
