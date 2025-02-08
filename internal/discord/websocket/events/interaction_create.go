package events

import (
	"slices"
	"strings"

	"presto/internal/bot/application_commands"
	"presto/internal/bot/message_components"
	"presto/internal/bot/modals"
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
	case api.INTERACTION_TYPE_MESSAGE_COMPONENT:
		HandleMessageComponentInteraction(interaction)
	}
}

func HandleApplicationCommands(interaction api.Interaction) {
	interactionName := discord.GetInteractionName(interaction.Data.Data)
	index := slices.IndexFunc(application_commands.RegisteredCommands, func(e application_commands.ApplicationCommandWithHandlers) bool {
		return discord.GetApplicationCommandName(e.Data) == interactionName
	})

	if index != -1 {
		var err error
		applicationCommand := application_commands.RegisteredCommands[index]

		if len(applicationCommand.Handlers) == 1 {
			err = applicationCommand.Handlers[0](interaction)
		} else {
			splittedInteractionName := strings.Split(interactionName, " ")
			applicationCommandToFind := splittedInteractionName[len(splittedInteractionName)-1]

			var subCommands []discord.ApplicationCommandOption
			var lookForSubCommands func([]discord.ApplicationCommandOption)

			lookForSubCommands = func(options []discord.ApplicationCommandOption) {
				for _, option := range options {
					if option.Type == discord.APPLICATION_COMMAND_OPTION_TYPE_SUB_COMMAND {
						subCommands = append(subCommands, option)
					} else if option.Type == discord.APPLICATION_COMMAND_OPTION_TYPE_SUB_COMMAND_GROUP {
						lookForSubCommands(option.Options)
					}
				}
			}

			lookForSubCommands(applicationCommand.Data.Options)

			for index, subCommand := range subCommands {
				if subCommand.Name == applicationCommandToFind {
					err = applicationCommand.Handlers[index](interaction)
				}
			}
		}

		if err != nil {
			interaction.RespondWithMessage(discord.Message{
				Embeds: []discord.Embed{
					{
						Description: err.Error(),
					},
				},
				Flags: discord.MESSAGE_FLAG_EPHEMERAL,
			})
		}
	}
}

func HandleModalSubmit(interaction api.Interaction) {
	registeredModals := modals.Get()
	index := slices.IndexFunc(registeredModals, func(e modals.WithHandler) bool {
		return e.Data.CustomID == interaction.Data.Data.CustomID
	})

	if index != -1 {
		err := registeredModals[index].Handler(interaction)
		modals.Remove(index)

		if err != nil {
			interaction.RespondWithMessage(discord.Message{
				Embeds: []discord.Embed{
					{
						Description: err.Error(),
					},
				},
				Flags: discord.MESSAGE_FLAG_EPHEMERAL,
			})
		}
	}
}

func HandleMessageComponentInteraction(interaction api.Interaction) {
	index := slices.IndexFunc(message_components.SelectMenus, func(e message_components.SelectMenuWithHandler) bool {
		return e.Data.CustomID == interaction.Data.Data.CustomID
	})

	if index != -1 {
		err := message_components.SelectMenus[index].Handler(interaction)
		if err != nil {
			interaction.RespondWithMessage(discord.Message{
				Embeds: []discord.Embed{
					{
						Description: err.Error(),
					},
				},
				Flags: discord.MESSAGE_FLAG_EPHEMERAL,
			})
			return
		}

		if message_components.SelectMenus[index].DeleteAfterInteracted {
			message_components.SelectMenus = append(message_components.SelectMenus[:index], message_components.SelectMenus[index+1:]...)
		}
	}
}
