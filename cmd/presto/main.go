package main

import (
	"log"
	"time"

	"presto/configs"
	"presto/internal/handlers/api"
	ws "presto/internal/handlers/websocket"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Started loading environment variables")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	configs.LoadDiscordConfig()

	log.Println("Finished loading environment variables successfully")

	log.Println("Started fetching Discord's websocket URL")
	gatewayData := api.GetGateway()
	log.Println("Finished fetching Discord's WSS URL successfully")

	log.Println("Started connection to Discord's WSS")
	connection := ws.ConnnectToWebsocket(gatewayData.URL)
	log.Println("Performed websocket handshake with Discord's WSS")

	log.Println("Started listening to WSS' messages")
	// This is the workaround I found in order to satisfy the data
	// that has been sent to Discord on every heartbeat.
	// Using `websocket.Conn.Ping()` method works "perfectly", but I
	// don't know if it's the best approach
	var lastEvent ws.DiscordGatewayEventPayload
	incomingEvents := make(chan ws.DiscordGatewayEventPayload)

	log.Println("Triggered initial handshake")
	ws.SendIdentify(connection)
	log.Println("Handshake ocurred successfully")

	go ws.ReadIncomingMessages(connection, incomingEvents)

	for {
		select {
		case <-time.After(40 * time.Second):
			log.Println("Started sending heartbeat")
			ws.SendHeartbeat(connection, lastEvent.SequenceNumber)
			log.Println("Sent heartbeat")
		case event := <-incomingEvents:
			lastEvent = event

			switch event.Name {
			case "READY":
				log.Println("Presto is ready to go")
			}
		}
	}
}
