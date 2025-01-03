package application_commands

import (
	"presto/internal/bot/message_components"
	"presto/internal/discord"
	"presto/internal/discord/api"
)

var Warn = NewUserCommand("warn", WarnHandler)

func WarnHandler(interaction api.Interaction) {
	modal := message_components.ModalWithHandler{
		Data: api.Modal{
			CustomID: "a",
			Title:    "Warning details",
			Components: []discord.MessageComponent{
				{
					Type: discord.MESSAGE_COMPONENT_TYPE_ACTION_ROW,
					Components: []discord.MessageComponent{
						{
							CustomID: "ab",
							Type:     discord.MESSAGE_COMPONENT_TYPE_TEXT_INPUT,
							Label:    "Reason for the warning",
							Style:    discord.TEXT_INPUT_STYLE_PARAGRAPH,
						},
					},
				},
			},
		},
		Handler: WarnModelHandler,
	}

	message_components.Modals = append(message_components.Modals, modal)
	interaction.RespondWithModal(modal.Data)
}

func WarnModelHandler(interaction api.Interaction) {
	interaction.RespondWithMessage(discord.Message{
		Content: "damn bro finally",
	})
}
