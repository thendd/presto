package application_commands

import (
	"fmt"
	"presto/internal/bot/errors"
	"presto/internal/bot/message_components"
	"presto/internal/bot/modals"
	"presto/internal/database"
	"presto/internal/discord"
	"presto/internal/discord/api"
	"presto/internal/log"
	"strconv"
	"time"
)

type SelectMenuInput struct {
	CommitChanges func(guild *database.Guild, value any)
	Conditions    func(value string) (any, error)
	TabName       string
}

var GuildSettings = NewSlashCommand(
	"settings",
	"Anything you want to customize",
	[]ApplicationCommandWithHandlerDataOption{},
	func(i api.Interaction) error { return nil },
).
	AddSubCommand("server", "Settings for your server", []ApplicationCommandWithHandlerDataOption{}, GuildSettingsHandler).
	ToApplicationCommand()

func GuildSettingsHandler(interaction api.Interaction) error {
	guild := database.Guild{
		ID: interaction.Data.GuildID,
	}

	if err := database.Connection.First(&guild).Error; err != nil {
		log.Error("There was an error when executing command \"settings\" invoked by the user %s at the guild %s when fetching the guild data: %s", interaction.Data.User.ID, interaction.Data.GuildID, err)
		return errors.UnknwonError
	}

	selectMenu := discord.MessageComponent{
		Type:        discord.MESSAGE_COMPONENT_TYPE_SELECT_MENU,
		CustomID:    strconv.Itoa(int(time.Now().UnixMilli())) + "-" + interaction.Data.GuildID + "-" + interaction.Data.Member.User.ID,
		Placeholder: "What do you want to configure?",
		Options: []discord.SelectOption{
			{
				Label: "Maximum quantity of warnings per user",
				Value: "0",
			},
			{
				Label:       "What happens when a user is warned too many times",
				Description: "Configure a punishment",
				Value:       "1",
			},
		},
	}

	switch guild.OnReachMaxWarningsPerUser {
	case int(database.ON_REACH_MAX_WARNINGS_PER_USER_BAN):
		selectMenu.Options = append(selectMenu.Options, discord.SelectOption{
			Label:       "Minutes to delete banned user messages for",
			Description: "If punishment is \"Ban\", deletes messages sent x minutes before",
			Value:       "2",
		})
	case int(database.ON_REACH_MAX_WARNINGS_PER_USER_GIVE_ROLE):
		selectMenu.Options = append(selectMenu.Options, discord.SelectOption{
			Label:       "Role to give to user",
			Description: "If punishment is \"Give role\", this role will be given to the user",
			Value:       "3",
		}, discord.SelectOption{
			Label:       "Minutes the user should keep the role for",
			Description: "If punishment is \"Give role\", decide how much time he will keep the role for",
			Value:       "4",
		})
	}

	message_components.SelectMenus = append(message_components.SelectMenus, message_components.SelectMenuWithHandler{
		Data:    selectMenu,
		Handler: ServerSettingsSelectMenuHandler,
		Args:    []any{guild, interaction.Data.Token},
	})

	err := interaction.RespondWithMessage(discord.Message{
		Components: []discord.MessageComponent{
			{
				Type: discord.MESSAGE_COMPONENT_TYPE_ACTION_ROW,
				Components: []discord.MessageComponent{
					selectMenu,
				},
			},
		},
		Flags: discord.MESSAGE_FLAG_EPHEMERAL,
	})
	if err != nil {
		log.Error(err.Error())
	}

	return nil
}

func ServerSettingsSelectMenuHandler(interaction api.Interaction, args ...any) error {
	guild := args[0].(database.Guild)
	originalInteractionToken := args[1].(string)
	settingsTab, _ := strconv.Atoi(interaction.Data.Data.Values[0])

	modalTemplate := api.Modal{
		CustomID: strconv.Itoa(int(time.Now().UnixMilli())) + "-" + interaction.Data.GuildID + "-" + interaction.Data.Member.User.ID,
		Components: []discord.MessageComponent{
			{
				Type: discord.MESSAGE_COMPONENT_TYPE_ACTION_ROW,
			},
		},
	}

	inputs := []SelectMenuInput{
		{
			TabName: "maximum amount of warnings a member can receive",
			Conditions: func(value string) (any, error) {
				newMaxWarningsPerUser, err := strconv.Atoi(value)

				if err != nil || newMaxWarningsPerUser < 1 {
					return 0, errors.New("Your answer must be a positive, whole number")
				}

				return newMaxWarningsPerUser, nil
			},
			CommitChanges: func(guild *database.Guild, value any) {
				guild.MaxWarningsPerUser = value.(int)
			},
		},
		{
			TabName: "punishment for a member that receives too many warnings",
			Conditions: func(value string) (any, error) {
				return strconv.Atoi(value)
			},
			CommitChanges: func(guild *database.Guild, value any) {
				guild.OnReachMaxWarningsPerUser = value.(int)
			},
		},
		{
			TabName: "quantity of minutes to delete banned member's messages for when they get too many warnings",
			Conditions: func(value string) (any, error) {
				newMinutesToDeleteUserMessagesFor, err := strconv.Atoi(value)

				if err != nil || newMinutesToDeleteUserMessagesFor < 1 || newMinutesToDeleteUserMessagesFor > 10080 {
					return 0, errors.New("Your answer must be a positive, whole number greater than 0 and lower than 10080")
				}

				return newMinutesToDeleteUserMessagesFor * 60, nil
			},
			CommitChanges: func(guild *database.Guild, value any) {
				guild.SecondsPunishedUserShouldKeepRoleFor = value.(int)
			},
		},
		{
			TabName: "role to give when a member gets too many warnings",
			Conditions: func(value string) (any, error) {
				return value, nil
			},
			CommitChanges: func(guild *database.Guild, value any) {
				guild.RoleToGiveOnReachMaxWarningsPerUser = value.(string)
			},
		},
		{
			TabName: "quantity of minutes the member should keep the role for when they get too many warnings",
			Conditions: func(value string) (any, error) {
				newMinutesUserShouldKeepRoleFor, err := strconv.Atoi(interaction.Data.Data.Components[0].Components[0].Value)
				if err != nil || newMinutesUserShouldKeepRoleFor < 1 {
					return 0, errors.New("Your answer must be a positive, whole number")
				}

				return newMinutesUserShouldKeepRoleFor * 60, nil
			},
			CommitChanges: func(guild *database.Guild, value any) {
				guild.SecondsPunishedUserShouldKeepRoleFor = value.(int)
			},
		},
	}

	var settingsTabHandler func(interaction api.Interaction, args ...any) error
	settingsTabHandler = func(interaction api.Interaction, _ ...any) error {
		guildToUpdate := database.Guild{
			ID: interaction.Data.GuildID,
		}

		var selectedTab string
		if len(interaction.Data.Data.Values) != 0 {
			selectedTab = interaction.Data.Data.Values[0]
		} else {
			selectedTab = interaction.Data.Data.Components[0].Components[0].Value
		}

		toCommit, err := inputs[settingsTab].Conditions(selectedTab)
		if err != nil {
			interaction.EditOriginalInteraction(discord.Message{
				Embeds: []discord.Embed{
					{
						Description: err.Error(),
						Color:       discord.EMBED_COLOR_RED,
					},
				},
				Flags: discord.MESSAGE_FLAG_EPHEMERAL,
			}, originalInteractionToken)

			return err
		}

		inputs[settingsTab].CommitChanges(&guildToUpdate, toCommit)

		if result := database.Connection.Save(guildToUpdate); result.Error != nil {
			return errors.UnknwonError
		}

		err = interaction.EditOriginalInteraction(discord.Message{
			Embeds: []discord.Embed{
				{
					Description: fmt.Sprintf("The **%s** was updated successfully.", inputs[settingsTab].TabName),
					Color:       discord.EMBED_COLOR_GREEN,
				},
			},
			Flags: discord.MESSAGE_FLAG_EPHEMERAL,
		}, originalInteractionToken)

		return nil
	}

	switch settingsTab {
	case 0:
		modalTemplate.Title = "Max warnings per user"
		modalTemplate.Components[0].Components = append(modalTemplate.Components[0].Components, discord.MessageComponent{
			CustomID:    "text-input",
			Type:        discord.MESSAGE_COMPONENT_TYPE_TEXT_INPUT,
			Style:       discord.TEXT_INPUT_STYLE_SHORT,
			Label:       "Your answer",
			Placeholder: "Ex: 3",
			Value:       strconv.Itoa(int(guild.MaxWarningsPerUser)),
			Required:    true,
			MinLength:   1,
			MaxLength:   2,
		})

		modals.Append(modals.WithHandler{
			Data:    modalTemplate,
			Handler: settingsTabHandler,
		})

		interaction.RespondWithModal(modalTemplate)
	case 1:
		selectMenu := discord.MessageComponent{
			Type:     discord.MESSAGE_COMPONENT_TYPE_SELECT_MENU,
			CustomID: strconv.Itoa(int(time.Now().UnixMilli())) + "-" + interaction.Data.GuildID + "-" + interaction.Data.ChannelID + "-" + interaction.Data.Member.User.ID,
			Options: []discord.SelectOption{
				{
					Label:       "Ban",
					Description: "Bans user. You can choose whether messages sent by them are deleted as well as the range of time",
					Value:       "0",
				},
				{
					Label:       "Kick",
					Description: "Simply kicks the user out of your server",
					Value:       "1",
				},
				{
					Label:       "Give role",
					Description: "You can choose which role and how much time the user should keep it",
					Value:       "2",
				},
			},
		}

		message_components.SelectMenus = append(message_components.SelectMenus, message_components.SelectMenuWithHandler{
			Data:    selectMenu,
			Handler: settingsTabHandler,
		})

		interaction.EditOriginalInteraction(discord.Message{
			Content: "What should I do when user is warned too many times?",
			Components: []discord.MessageComponent{
				{
					Type:       discord.MESSAGE_COMPONENT_TYPE_ACTION_ROW,
					Components: []discord.MessageComponent{selectMenu},
				},
			},
			Flags: discord.MESSAGE_FLAG_EPHEMERAL,
		}, originalInteractionToken)
	case 2:
		modalTemplate.Title = "Quantity of minutes to delete banned user messages for"
		modalTemplate.Components[0].Components = append(modalTemplate.Components[0].Components, discord.MessageComponent{
			CustomID:    "text-input",
			Type:        discord.MESSAGE_COMPONENT_TYPE_TEXT_INPUT,
			Style:       discord.TEXT_INPUT_STYLE_SHORT,
			Label:       "Your answer (0 is no messages)",
			Placeholder: "Ex: 10",
			Value:       strconv.Itoa(int(guild.SecondsToDeleteMessagesForOnReachMaxWarningsPerUser / 60)),
			Required:    true,
			MinLength:   1,
			MaxLength:   5,
		})

		modals.Append(modals.WithHandler{
			Data:    modalTemplate,
			Handler: settingsTabHandler,
		})

		interaction.RespondWithModal(modalTemplate)
	case 3:
		selectMenu := discord.MessageComponent{
			Type:     discord.MESSAGE_COMPONENT_TYPE_ROLE_SELECT,
			CustomID: strconv.Itoa(int(time.Now().UnixMilli())) + "-" + interaction.Data.GuildID + "-" + interaction.Data.ChannelID + "-" + interaction.Data.Member.User.ID,
		}

		message_components.SelectMenus = append(message_components.SelectMenus, message_components.SelectMenuWithHandler{
			Data:    selectMenu,
			Handler: settingsTabHandler,
		})

		interaction.EditOriginalInteraction(discord.Message{
			Content: "What role should be given to the user that is warned too many times?",
			Components: []discord.MessageComponent{
				{
					Type:       discord.MESSAGE_COMPONENT_TYPE_ACTION_ROW,
					Components: []discord.MessageComponent{selectMenu},
				},
			},
			Flags: discord.MESSAGE_FLAG_EPHEMERAL,
		}, originalInteractionToken)
	case 4:
		modalTemplate.Title = "Minutes the user should keep the role for"
		modalTemplate.Components[0].Components = append(modalTemplate.Components[0].Components, discord.MessageComponent{
			CustomID:    "text-input",
			Type:        discord.MESSAGE_COMPONENT_TYPE_TEXT_INPUT,
			Style:       discord.TEXT_INPUT_STYLE_SHORT,
			Label:       "Your answer",
			Placeholder: "Ex: 10",
			Value:       strconv.Itoa(int(guild.SecondsPunishedUserShouldKeepRoleFor / 60)),
			Required:    true,
			MinLength:   1,
			MaxLength:   5,
		})

		modals.Append(modals.WithHandler{
			Data:    modalTemplate,
			Handler: settingsTabHandler,
		})

		interaction.RespondWithModal(modalTemplate)
	}

	return nil
}
