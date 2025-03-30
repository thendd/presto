package commands

import (
	"fmt"
	"presto/internal/bot"
	"presto/internal/bot/errors"
	"presto/internal/log"

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
	args := []any{context.Interaction.Data.Data.Options[0].Value.(string)}

	switch context.Interaction.Data.Data.Type {
	case discord.APPLICATION_COMMAND_TYPE_USER:
		args[0] = context.Interaction.Data.Data.TargetID
	case discord.APPLICATION_COMMAND_TYPE_MESSAGE:
		message := context.Interaction.Data.Data.Resolved.Messages[context.Interaction.Data.Data.TargetID]
		args[0] = message.Author.ID
		args = append(args, message.ChannelID, message.ID)
		textInputLabel = "Additional info"
		isTextInputRequired = false
	}

	modal := bot.NewModalWithHandler(discord.Modal{
		Title: "Warning details",
		Components: []discord.MessageComponent{
			{
				Type: discord.MESSAGE_COMPONENT_TYPE_ACTION_ROW,
				Components: []discord.MessageComponent{
					discord.NewTextInputComponent(discord.MessageComponent{

						Label:     textInputLabel,
						Style:     discord.TEXT_INPUT_STYLE_PARAGRAPH,
						MaxLength: 1500,
						Required:  isTextInputRequired,
					}),
				},
			},
		},
	}, WarnModalHandler, []any{})

	context.Session.Cache.Modals.Append(modal)
	context.Interaction.RespondWithModal(modal.Data)

	return nil
}

func WarnModalHandler(context bot.Context, args ...any) error {
	targetId := args[0].(string)

	guildData := &database.Guild{
		ID: context.Interaction.Data.GuildID,
	}

	if result := database.Connection.First(&guildData); result.Error != nil {
		log.Error("There was an error when executing command \"warn\" invoked by the user %s at the guild %s when fetching the guild data: %s", context.Interaction.Data.User.ID, context.Interaction.Data.GuildID, result.Error)
		return errors.UnknwonError
	}

	target := database.GuildMember{
		GuildId: context.Interaction.Data.GuildID,
		UserId:  targetId,
	}
	if result := database.Connection.FirstOrCreate(&target); result.Error != nil {
		log.Error("There was an error when executing command \"warn\" invoked by the user %s at the guild %s when fetching the target data: %s", context.Interaction.Data.User.ID, context.Interaction.Data.GuildID, result.Error)
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

	if len(args) > 1 {
		warningEmbedDescription += fmt.Sprintf("**[This message](%s)** was attached to the warning.", fmt.Sprintf("https://discord.com/channels/%s/%s/%s", context.Interaction.Data.GuildID, args[1], args[2]))
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
