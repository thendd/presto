package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"presto/internal/bot/application_commands"
	"presto/internal/config"
	"presto/internal/database"
	"presto/internal/discord"
	"presto/internal/discord/api"
	ws "presto/internal/discord/websocket"
	"presto/internal/discord/websocket/events"
)

func main() {
	config.LoadEnvironmentVariables()
	database.Connect()
	defer database.Connection.Close(context.Background())

	gatewayData := api.GetGateway()
	ws.Connnect(gatewayData.URL)

	// This is the workaround I found in order to satisfy the data
	// that has to be sent to Discord on every heartbeat.
	// Using `websocket.Conn.Ping()` method works "perfectly", but I
	// don't know if it's the best approach
	var lastEvent ws.EventPayload
	incomingEvents := make(chan ws.EventPayload)

	events.SendIdentify()
	application_commands.Register()
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
			case events.GUILD_CREATE:
				// Even though an unavailable guild object might be sent,
				// it would still have the id property, which is the only one
				// that will be used
				var guildData discord.Guild
				json.Unmarshal(eventData, &guildData)

				queries := database.New(database.Connection)
				createdGuild, _ := queries.CreateGuild(context.Background(), guildData.ID)
				if createdGuild.ID != "" {
					log.Printf("Created guild (%s) in database successfully\n", guildData.ID)
				}
			}
		}
	}
}
