package ws

import (
	"slices"

	"presto/internal/bot/application_commands"
	"presto/internal/bot/message_components"
	"presto/internal/bot/modals"
	"presto/internal/discord"
	"presto/internal/discord/api"
)

func ReceiveInteractionCreate(interactionData discord.InteractionCreatePayload) {
	interaction := api.Interaction{
		Data: interactionData,
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
	applicationCommandIndex := slices.IndexFunc(application_commands.Local, func(e application_commands.ApplicationCommandWithHandler) bool {
		return slices.Contains(discord.GetFullNamesOfApplicationCommand(e.ToApplicationCommand()), interactionName)
	})

	if applicationCommandIndex != -1 {
		var err error

		applicationCommand := application_commands.Local[applicationCommandIndex]
		if applicationCommand.Data.Name == interactionName {
			err = applicationCommand.Handler(interaction)
		} else {
			for _, option := range applicationCommand.Data.Options {
				if applicationCommand.Data.Name+" "+option.Name == interactionName {
					err = option.Handler(interaction)
					break
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
		err := registeredModals[index].Handler(interaction, registeredModals[index].Args...)
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
		err := message_components.SelectMenus[index].Handler(interaction, message_components.SelectMenus[index].Args...)
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

		message_components.SelectMenus = slices.Delete(message_components.SelectMenus, index, index+1)
	}
}
