package events

import (
	"slices"

	"presto/internal/bot"
	"presto/internal/bot/application_commands"
	"presto/internal/bot/message_components"
	"presto/internal/discord"
	"presto/internal/discord/api"
	ws "presto/internal/discord/websocket"
)

func ReceiveInteractionCreate(interactionData discord.InteractionCreatePayload) {
	interaction := api.Interaction{
		Data:      interactionData,
		Websocket: ws.Connection,
	}

	switch interactionData.Type {
	case api.INTERACTION_TYPE_APPLICATION_COMMAND:
		HandleApplicationCommands(interaction)
	case api.INTERACTION_TYPE_MODAL_SUBMIT:
		HandleModalSubmit(interaction)
	}
}

func HandleApplicationCommands(interaction api.Interaction) {
	index := slices.IndexFunc(bot.RegisteredCommands, func(e application_commands.ApplicationCommandWithHandler) bool {
		return e.Data.Name == interaction.Data.Data.Name
	})

	if index != -1 {
		bot.RegisteredCommands[index].Handler(interaction)
	}
}

func HandleModalSubmit(interaction api.Interaction) {
	index := slices.IndexFunc(message_components.Modals, func(e message_components.ModalWithHandler) bool {
		return e.Data.CustomID == interaction.Data.Data.CustomID
	})

	if index != -1 {
		message_components.Modals[index].Handler(interaction)
		message_components.Modals = append(message_components.Modals[:index], message_components.Modals[index+1:]...)
	}
}
