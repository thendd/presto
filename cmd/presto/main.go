package main

import (
	"encoding/json"
	"log"
	"time"

	"presto/internal/bot"
	"presto/internal/config"
	"presto/internal/discord"
	"presto/internal/discord/api"
	ws "presto/internal/discord/websocket"
	"presto/internal/discord/websocket/events"
)

func main() {
	config.LoadEnvironmentVariables()
	gatewayData := api.GetGateway()
	ws.Connnect(gatewayData.URL)

	// This is the workaround I found in order to satisfy the data
	// that has to be sent to Discord on every heartbeat.
	// Using `websocket.Conn.Ping()` method works "perfectly", but I
	// don't know if it's the best approach
	var lastEvent ws.EventPayload
	incomingEvents := make(chan ws.EventPayload)

	events.SendIdentify()
	bot.RegisterCommands()
	go ws.OnEventReceive(incomingEvents)

	for {
		select {
		case <-time.After(40 * time.Second):
			events.SendHeartbeat(lastEvent.SequenceNumber)
		case event := <-incomingEvents:
			lastEvent = event

			eventData, _ := json.Marshal(event.Data)

			switch event.Name {
			case events.READY:
				log.Println("Presto is ready to go")
			case events.INTERACTION_CREATE:
				var interactionData discord.InteractionCreatePayload
				json.Unmarshal(eventData, &interactionData)

				events.ReceiveInteractionCreate(interactionData)
			}
		}
	}
}
