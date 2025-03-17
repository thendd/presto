package commands

import (
	"fmt"
	"presto/internal/bot"
	"presto/internal/bot/errors"
	"presto/internal/log"
	"strings"
	"time"

	"presto/internal/database"
	"presto/internal/discord"
)

var (
	WarnUserCommand  = bot.NewUserCommand("Warn", WarnHandler)
	WarnSlashCommand = bot.NewSlashCommand("warn", "Sends a warning to the user", []bot.ApplicationCommandWithHandlerDataOption{
		{
			Type:        discord.APPLICATION_COMMAND_OPTION_TYPE_USER,
			Name:        "user",
			Description: "The user you would like to send a warning to",
			Required:    true,
		},
	}, WarnHandler).ToApplicationCommand()
	WarnMessageCommand = bot.NewMessageCommand("Warn", WarnHandler)
)

func WarnHandler(context bot.Context) error {
	textInputLabel := "Reason"
	isTextInputRequired := true
	modalCustomId := fmt.Sprintf("%d-%s-", time.Now().UnixMilli(), context.Interaction.Data.Member.User.ID)

	switch context.Interaction.Data.Data.Type {
	case discord.APPLICATION_COMMAND_TYPE_CHAT_INPUT:
		modalCustomId += context.Interaction.Data.Data.Options[0].Value.(string)
	case discord.APPLICATION_COMMAND_TYPE_USER:
		modalCustomId += context.Interaction.Data.Data.TargetID
	case discord.APPLICATION_COMMAND_TYPE_MESSAGE:
		message := context.Interaction.Data.Data.Resolved.Messages[context.Interaction.Data.Data.TargetID]
		modalCustomId += message.Author.ID + "-" + message.ChannelID + "-" + message.ID
		textInputLabel = "Additional info"
		isTextInputRequired = false
	}

	modal := bot.ModalWithHandler{
		Data: discord.Modal{
			CustomID: modalCustomId,
			Title:    "Warning details",
			Components: []discord.MessageComponent{
				{
					Type: discord.MESSAGE_COMPONENT_TYPE_ACTION_ROW,
					Components: []discord.MessageComponent{
						{
							CustomID:  modalCustomId + "-0",
							Type:      discord.MESSAGE_COMPONENT_TYPE_TEXT_INPUT,
							Label:     textInputLabel,
							Style:     discord.TEXT_INPUT_STYLE_PARAGRAPH,
							MaxLength: 1500,
							Required:  isTextInputRequired,
						},
					},
				},
			},
		},
		Handler: WarnModelHandler,
	}

	context.Session.Cache.Modals.Append(modal)
	context.Interaction.RespondWithModal(modal.Data)

	return nil
}

func WarnModelHandler(context bot.Context, _ ...any) error {
	splittedCustomID := strings.Split(context.Interaction.Data.Data.CustomID, "-")

	targetId := strings.Split(context.Interaction.Data.Data.CustomID, "-")[2]

	guildData := &database.Guild{
		ID: context.Interaction.Data.GuildID,
	}

	if result := database.Connection.First(&guildData); result.Error != nil {
		log.Errorf("There was an error when executing command \"warn\" invoked by the user %s at the guild %s when fetching the guild data: %s", context.Interaction.Data.User.ID, context.Interaction.Data.GuildID, result.Error)
		return errors.UnknwonError
	}

	target := database.GuildMember{
		GuildId: context.Interaction.Data.GuildID,
		UserId:  targetId,
	}
	if result := database.Connection.FirstOrCreate(&target); result.Error != nil {
		log.Errorf("There was an error when executing command \"warn\" invoked by the user %s at the guild %s when fetching the target data: %s", context.Interaction.Data.User.ID, context.Interaction.Data.GuildID, result.Error)
		return errors.UnknwonError
	}

	remainingWarnings := guildData.MaxWarningsPerUser - target.Warnings - 1

	dmChannel, err := context.Session.Cache.DMChannels.GetByRecipientID(targetId)
	if err != nil {
		dmChannel = discord.CreateDM(targetId)
	}

	guild, err := context.Session.Cache.Guilds.GetByID(context.Interaction.Data.GuildID)
	if err != nil {
		guild = discord.GetGuildById(context.Interaction.Data.GuildID)
	}

	if remainingWarnings < 0 {
		context.Interaction.RespondWithMessage(discord.Message{
			Embeds: []discord.Embed{
				{
					Description: "As the user has already reached the limit of warnings, they will be banned.",
					Color:       discord.EMBED_COLOR_GREEN,
				},
			},
			Flags: discord.MESSAGE_FLAG_EPHEMERAL,
		})

		discord.BanMember(guildData.ID, targetId)

		message := discord.Message{
			ChannelID: dmChannel.ID,
			Embeds: []discord.Embed{
				{
					Description: fmt.Sprintf("You were banned from %s because you have received too many warnings.", guild.Name),
					Color:       discord.EMBED_COLOR_RED,
				},
			},
		}

		message.Send()

		return nil
	}

	context.Interaction.RespondWithMessage(discord.Message{
		Embeds: []discord.Embed{
			{
				Description: fmt.Sprintf("Warning was sent successfully. The user will receive a message in their DM with the reason of the warning. They still have **%d warnings left**.", remainingWarnings),
				Color:       discord.EMBED_COLOR_GREEN,
			},
		},
		Flags: discord.MESSAGE_FLAG_EPHEMERAL,
	})

	database.Connection.Save(&database.GuildMember{
		UserId:   targetId,
		GuildId:  context.Interaction.Data.GuildID,
		Warnings: target.Warnings + 1,
	})

	warningEmbedDescription := fmt.Sprintf("You were warned in the server **%s** ", guild.Name)
	if context.Interaction.Data.Data.Components[0].Components[0].Value != "" {
		warningEmbedDescription += fmt.Sprintf("with the following reason: **%s**.", context.Interaction.Data.Data.Components[0].Components[0].Value)
	} else {
		warningEmbedDescription += "but no reason was given."
	}

	if len(splittedCustomID) == 5 {
		warningEmbedDescription += fmt.Sprintf("**[This message](%s)** was attached to the warning.", fmt.Sprintf("https://discord.com/channels/%s/%s/%s", context.Interaction.Data.GuildID, splittedCustomID[3], splittedCustomID[4]))
	}

	warningEmbed := discord.Embed{
		Description: warningEmbedDescription,
		Color:       discord.EMBED_COLOR_YELLOW,
		Image: &discord.EmbedImage{
			URL:    context.Interaction.Data.Guild.GetIconURL(),
			Height: 100,
			Width:  100,
		},
	}

	message := discord.Message{
		ChannelID: dmChannel.ID,
		Embeds:    []discord.Embed{warningEmbed},
	}
	message.Send()

	return nil
}
