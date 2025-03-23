package bot

import (
	"slices"

	"presto/internal/discord"
)

func (session *Session) HandleInteractionCreateEvent(interactionData discord.InteractionCreatePayload) {
	interaction := discord.Interaction{
		Data: interactionData,
	}

	switch interactionData.Type {
	case discord.INTERACTION_TYPE_APPLICATION_COMMAND:
		handleApplicationCommands(session, interaction)
	case discord.INTERACTION_TYPE_MODAL_SUBMIT:
		handleModalSubmit(session, interaction)
	case discord.INTERACTION_TYPE_MESSAGE_COMPONENT:
		handleMessageComponentInteraction(session, interaction)
	}
}

func handleApplicationCommands(session *Session, interaction discord.Interaction) {
	localApplicationCommands := session.RegisteredCommands

	interactionName := discord.GetInteractionName(interaction.Data.Data)
	applicationCommandIndex := slices.IndexFunc(localApplicationCommands, func(e ApplicationCommandWithHandler) bool {
		return slices.Contains(discord.GetFullNamesOfApplicationCommand(e.ToApplicationCommand()), interactionName)
	})

	if applicationCommandIndex != -1 {
		var err error

		applicationCommand := localApplicationCommands[applicationCommandIndex]
		if applicationCommand.Data.Name == interactionName {
			err = applicationCommand.Handler(Context{
				Session:     session,
				Interaction: interaction,
			})
		} else {
			for _, option := range applicationCommand.Data.Options {
				if applicationCommand.Data.Name+" "+option.Name == interactionName {
					err = option.Handler(Context{
						Session:     session,
						Interaction: interaction,
					})
					break
				}
			}
		}

		if err != nil {
			interaction.RespondWithMessage(discord.Message{
				Embeds: []discord.Embed{
					{
						Description: err.Error(),
						Color:       discord.EMBED_COLOR_RED,
					},
				},
				Flags: discord.MESSAGE_FLAG_EPHEMERAL,
			})
		}
	}
}

func handleModalSubmit(session *Session, interaction discord.Interaction) {
	index := slices.IndexFunc(session.Cache.Modals, func(e ModalWithHandler) bool {
		return e.Data.CustomID == interaction.Data.Data.CustomID
	})

	if index != -1 {
		err := session.Cache.Modals[index].Handler(Context{
			Session:     session,
			Interaction: interaction,
		}, session.Cache.Modals[index].Args...)
		session.Cache.Modals.Remove(index)

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

func handleMessageComponentInteraction(session *Session, interaction discord.Interaction) {
	index := slices.IndexFunc(session.Cache.SelectMenus, func(e SelectMenuWithHandler) bool {
		return e.Data.CustomID == interaction.Data.Data.CustomID
	})

	if index != -1 {
		err := session.Cache.SelectMenus[index].Handler(Context{
			Session:     session,
			Interaction: interaction,
		}, session.Cache.SelectMenus[index].Args...)
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

		session.Cache.SelectMenus.Remove(index)

		return
	}

	interaction.RespondWithMessage(discord.Message{
		Embeds: []discord.Embed{
			{
				Description: "There was an unexpected error while processing your interaction.",
				Color:       discord.EMBED_COLOR_RED,
			},
		},
		Flags: discord.MESSAGE_FLAG_EPHEMERAL,
	})
}
