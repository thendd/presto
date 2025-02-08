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

var Settings = NewSlashCommandGroup("settings", "Anything you want to customize").
	AddSubCommandGroup("server", "Settings for your server").
	AddSubCommand("server", "warnings", "What should I do when a user is warned?", []discord.ApplicationCommandOption{}, ServerWarningSettingsHandler).
	ToApplicationCommand()

func ServerWarningSettingsHandler(interaction api.Interaction) error {
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
			{
				Label:       "Minutes to delete banned user messages for",
				Description: "If punishment is \"Ban\", deletes messages from the user x minutes before",
				Value:       "2",
			},
			{
				Label:       "Role to give to user",
				Description: "If punishment is \"Give role\", this role will be given to the user",
				Value:       "3",
			},
			{
				Label:       "Minutes the user should keep the role for",
				Description: "If punishment is \"Give role\", decide how much time he will keep the role for",
				Value:       "4",
			},
		},
	}

	message_components.SelectMenus = append(message_components.SelectMenus, message_components.SelectMenuWithHandler{
		Data:                  selectMenu,
		Handler:               ServerWarningSettingsSelectMenuHandler,
		DeleteAfterInteracted: false,
	})

	interaction.RespondWithMessage(discord.Message{
		Components: []discord.MessageComponent{
			{
				Type: discord.MESSAGE_COMPONENT_TYPE_ACTION_ROW,
				Components: []discord.MessageComponent{
					selectMenu,
				},
			},
		},
	})

	return nil
}

func ServerWarningSettingsSelectMenuHandler(interaction api.Interaction) error {
	guild := database.Guild{
		ID: interaction.Data.GuildID,
	}
	result := database.Connection.First(guild)
	if result.Error != nil {
		log.Error("There was an error when executing command \"settings\" invoked by the user %s at the guild %s when fetching the server data: %s", interaction.Data.User.ID, interaction.Data.GuildID, result.Error)
		return errors.UnknwonError
	}

	currentOnReachMaxWarningsPerUser := guild.OnReachMaxWarningsPerUser

	settingsTab, _ := strconv.Atoi(interaction.Data.Data.Values[0])

	template := "This tab is only accessible if the punishment for a user that gets too many warnings is **%s**"

	if settingsTab == 2 && currentOnReachMaxWarningsPerUser != int8(database.ON_REACH_MAX_WARNINGS_PER_USER_BAN) {
		return errors.New(fmt.Sprintf(template, "Ban user"))
	} else if (settingsTab == 3 || settingsTab == 4) && currentOnReachMaxWarningsPerUser != int8(database.ON_REACH_MAX_WARNINGS_PER_USER_GIVE_ROLE) {
		return errors.New(fmt.Sprintf(template, "Give role"))
	}

	modalTemplate := api.Modal{
		CustomID: strconv.Itoa(int(time.Now().UnixMilli())) + "-" + interaction.Data.GuildID + "-" + interaction.Data.Member.User.ID,
		Components: []discord.MessageComponent{
			{
				Type: discord.MESSAGE_COMPONENT_TYPE_ACTION_ROW,
			},
		},
	}

	var settingsTabHandler func(interaction api.Interaction) error
	settingsTabHandler = func(interaction api.Interaction) error {
		successResponse := discord.Message{
			Embeds: []discord.Embed{
				{
					Description: "The **%s** was updated successfully.",
					Color:       discord.EMBED_COLOR_GREEN,
				},
			},
			Flags: discord.MESSAGE_FLAG_EPHEMERAL,
		}

		guildToUpdate := database.Guild{
			ID: interaction.Data.GuildID,
		}

		switch settingsTab {
		case 0:
			newMaxWarningsPerUser, err := strconv.Atoi(interaction.Data.Data.Components[0].Components[0].Value)

			if err != nil || newMaxWarningsPerUser < 0 {
				return errors.New("Your answer must be a positive, whole number or zero")
			}

			guildToUpdate.MaxWarningsPerUser = int8(newMaxWarningsPerUser)

			successResponse.Embeds[0].Description = fmt.Sprintf(successResponse.Embeds[0].Description, "maximum amount of warnings a user can receive")
		case 1:
			newOnReachMaxWarningsPerUser, _ := strconv.Atoi(interaction.Data.Data.Values[0])

			guildToUpdate.OnReachMaxWarningsPerUser = int8(newOnReachMaxWarningsPerUser)

			successResponse.Embeds[0].Description = fmt.Sprintf(successResponse.Embeds[0].Description, "punishment for a user that receives too many warnings")
		case 2:
			newMinutesToDeleteUserMessagesFor, err := strconv.Atoi(interaction.Data.Data.Components[0].Components[0].Value)

			if err != nil || newMinutesToDeleteUserMessagesFor < 1 || newMinutesToDeleteUserMessagesFor > 10080 {
				return errors.New("Your answer must be a positive, whole number greater than 0 and lower than 10080")
			}

			guildToUpdate.SecondsToDeleteMessagesForOnReachMaxWarningsPerUser = newMinutesToDeleteUserMessagesFor * 60

			successResponse.Embeds[0].Description = fmt.Sprintf(successResponse.Embeds[0].Description, "quantity of minutes to delete banned user's messages for when they get too many warnings")
		case 3:
			newRoleToGiveOnReachMaxWarningsPerUser := interaction.Data.Data.Values[0]

			guildToUpdate.RoleToGiveOnReachMaxWarningsPerUser = newRoleToGiveOnReachMaxWarningsPerUser

			successResponse.Embeds[0].Description = fmt.Sprintf(successResponse.Embeds[0].Description, "role to give when the user is warned too many times")
		case 4:
			newMinutesUserShouldKeepRoleFor, err := strconv.Atoi(interaction.Data.Data.Components[0].Components[0].Value)
			if err != nil || newMinutesUserShouldKeepRoleFor < 0 {
				return errors.New("Your answer must be a positive, whole number")
			}

			guildToUpdate.SecondsPunishedUserShouldKeepRoleFor = newMinutesUserShouldKeepRoleFor * 60

			successResponse.Embeds[0].Description = fmt.Sprintf(successResponse.Embeds[0].Description, "quantity of minutes the user should keep the role for when they get too many warnings")
		}

		if result := database.Connection.Save(guildToUpdate); result.Error != nil {
			return errors.UnknwonError
		}

		interaction.RespondWithMessage(successResponse)

		return nil
	}

	switch settingsTab {
	case 0:
		modalTemplate.Title = "Max warnings per user"
		modalTemplate.Components[0].Components = append(modalTemplate.Components[0].Components, discord.MessageComponent{
			CustomID:    "text-input",
			Type:        discord.MESSAGE_COMPONENT_TYPE_TEXT_INPUT,
			Style:       discord.TEXT_INPUT_STYLE_SHORT,
			Label:       "Your answer (0 is unlimited)",
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
				{
					Label: "Nothing",
					Value: "3",
				},
			},
		}

		message_components.SelectMenus = append(message_components.SelectMenus, message_components.SelectMenuWithHandler{
			Data:                  selectMenu,
			DeleteAfterInteracted: true,
			Handler:               settingsTabHandler,
		})

		interaction.RespondWithMessage(discord.Message{
			Content: "What should I do when user is warned too many times?",
			Components: []discord.MessageComponent{
				{
					Type:       discord.MESSAGE_COMPONENT_TYPE_ACTION_ROW,
					Components: []discord.MessageComponent{selectMenu},
				},
			},
			Flags: discord.MESSAGE_FLAG_EPHEMERAL,
		})
	case 2:
		modalTemplate.Title = "Minutes to delete banned user messages for"
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
			Data:                  selectMenu,
			DeleteAfterInteracted: true,
			Handler:               settingsTabHandler,
		})

		interaction.RespondWithMessage(discord.Message{
			Content: "What role should be given to the user that is warned too many times?",
			Components: []discord.MessageComponent{
				{
					Type:       discord.MESSAGE_COMPONENT_TYPE_ACTION_ROW,
					Components: []discord.MessageComponent{selectMenu},
				},
			},
			Flags: discord.MESSAGE_FLAG_EPHEMERAL,
		})
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
