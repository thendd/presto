package main

import (
	"encoding/json"
	"errors"
	"presto/internal/log"
	"time"

	"presto/internal/bot/application_commands"
	"presto/internal/config"
	"presto/internal/database"
	"presto/internal/discord"
	"presto/internal/discord/api"
	ws "presto/internal/discord/websocket"
	"presto/internal/discord/websocket/events"

	"gorm.io/gorm"
)

func main() {
	config.LoadEnvironmentVariables()
	database.Connect()

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
				log.Info("Presto is ready to go")
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

				result := database.Connection.Create(&database.Guild{
					ID: guildData.ID,
				})
				if result.Error != nil && !errors.Is(result.Error, gorm.ErrDuplicatedKey) {
					log.Error("Failed to create guild %s: %s", guildData.ID, result.Error)
					continue
				} else if result.Error != nil {
					log.Error("Guild %s was not created because it already exists in the database", guildData.ID)
					continue
				}

				log.Info("Created guild (%s) in database successfully\n", guildData.ID)
			}
		}
	}
}
