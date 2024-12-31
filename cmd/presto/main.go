package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"presto/configs"
	"presto/internal/handlers/api"
	ws "presto/internal/handlers/websocket"

	"github.com/coder/websocket"
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
	var lastMessage map[string]any
	incomingMessages := make(chan map[string]any)

	log.Println("Triggered initial handshake")
	ws.Identify(connection)
	log.Println("Handshake ocurred successfully")

	go ws.ReadIncomingMessages(connection, incomingMessages)

	for {
		select {
		case <-time.After(40 * time.Second):
			log.Println("Started sending heartbeat")
			var heartbeat string

			lastSequenceNumber := lastMessage["s"]
			if lastSequenceNumber != nil {
				heartbeat = fmt.Sprintf("{\"op\": 1, \"d\": %v}", lastSequenceNumber)
			} else {
				heartbeat = "{\"op\": 1, \"d\": null}"
			}

			err := connection.Write(context.Background(), websocket.MessageText, []byte(heartbeat))
			if err != nil {
				log.Fatal(err)
			}

			log.Println("Sent heartbeat")
		case message := <-incomingMessages:
			lastMessage = message
			if message["t"] == "READY" {
				log.Println("Presto is ready to go")
			}
		}
	}
}
