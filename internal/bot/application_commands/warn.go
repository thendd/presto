package application_commands

import (
	"fmt"
	"log"
	"strings"
	"time"

	"presto/internal/bot/modals"
	"presto/internal/database"
	"presto/internal/discord"
	"presto/internal/discord/api"
	"presto/internal/discord/api/cache"
	"presto/internal/discord/cdn"
)

var (
	WarnUserCommand  = NewUserCommand("Warn", WarnHandler)
	WarnSlashCommand = NewSlashCommand("warn", "Sends a warning to the user", []discord.ApplicationCommandOption{
		{
			Type:        discord.APPLICATION_COMMAND_OPTION_TYPE_USER,
			Name:        "user",
			Description: "The user you would like to send a warning to",
			Required:    true,
		},
	}, WarnHandler)
	WarnMessageCommand = NewMessageCommand("Warn", WarnHandler)
)

func WarnHandler(interaction api.Interaction) {
	textInputLabel := "Reason"
	isTextInputRequired := true
	modalCustomId := fmt.Sprintf("%d-%s-", time.Now().UnixMilli(), interaction.Data.Member.User.ID)

	switch interaction.Data.Data.Type {
	case discord.APPLICATION_COMMAND_TYPE_CHAT_INPUT:
		modalCustomId += interaction.Data.Data.Options[0].Value.(string)
	case discord.APPLICATION_COMMAND_TYPE_USER:
		modalCustomId += interaction.Data.Data.TargetID
	case discord.APPLICATION_COMMAND_TYPE_MESSAGE:
		message := interaction.Data.Data.Resolved.Messages[interaction.Data.Data.TargetID]
		modalCustomId += message.Author.ID + "-" + message.ChannelID + "-" + message.ID
		textInputLabel = "Additional info"
		isTextInputRequired = false
	}

	modal := modals.WithHandler{
		Data: api.Modal{
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

	modals.Append(modal)
	interaction.RespondWithModal(modal.Data)
}

func WarnModelHandler(interaction api.Interaction) {
	splittedCustomID := strings.Split(interaction.Data.Data.CustomID, "-")

	targetId := strings.Split(interaction.Data.Data.CustomID, "-")[2]

	guildData := &database.Guild{
		ID: interaction.Data.GuildID,
	}

	result := database.Connection.First(&guildData)
	if result.Error != nil {
		log.Printf("There was an error when executing command \"warn\" invoked by the user %s at the guild %s when fetching the guild data: %s", interaction.Data.User.ID, interaction.Data.GuildID, result.Error)
		interaction.RespondWithMessage(discord.Message{
			Embeds: []discord.Embed{
				{
					Description: "This command did not work due to an unknown error.",
				},
			},
			Flags: discord.MESSAGE_FLAG_EPHEMERAL,
		})
		return
	}

	target := database.GuildMember{
		GuildId: interaction.Data.GuildID,
		UserId:  targetId,
	}
	result = database.Connection.FirstOrCreate(&target)
	if result.Error != nil {
		log.Printf("There was an error when executing command \"warn\" invoked by the user %s at the guild %s when fetching the target data: %s", interaction.Data.User.ID, interaction.Data.GuildID, result.Error)
		interaction.RespondWithMessage(discord.Message{
			Embeds: []discord.Embed{
				{
					Description: "This command did not work due to an unknown error.",
				},
			},
			Flags: discord.MESSAGE_FLAG_EPHEMERAL,
		})
		return
	}

	remainingWarnings := guildData.MaxWarningsPerUser - target.Warnings - 1

	dmChannel := cache.GetDMChannelByRecipientID(targetId)
	if dmChannel.ID == "" {
		dmChannel = api.CreateDM(targetId)
	}

	guild := cache.GetGuildById(interaction.Data.GuildID)
	if guild.Name == "" {
		guild = api.GetGuildById(interaction.Data.GuildID)
	}

	if remainingWarnings < 0 {
		interaction.RespondWithMessage(discord.Message{
			Embeds: []discord.Embed{
				{
					Description: "As the user has already reached the limit of warnings, they will be banned.",
					Color:       discord.EMBED_COLOR_GREEN,
				},
			},
			Flags: discord.MESSAGE_FLAG_EPHEMERAL,
		})

		api.BanUser(guildData.ID, targetId)

		api.SendMessage(discord.Message{
			ChannelID: dmChannel.ID,
			Embeds: []discord.Embed{
				{
					Description: fmt.Sprintf("You were banned from %s because you have received too many warnings.", guild.Name),
					Color:       discord.EMBED_COLOR_RED,
				},
			},
		})

		return
	}

	interaction.RespondWithMessage(discord.Message{
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
		GuildId:  interaction.Data.GuildID,
		Warnings: target.Warnings + 1,
	})

	warningEmbedDescription := fmt.Sprintf("You were warned in the server **%s** ", guild.Name)
	if interaction.Data.Data.Components[0].Components[0].Value != "" {
		warningEmbedDescription += fmt.Sprintf("with the following reason: **%s**.", interaction.Data.Data.Components[0].Components[0].Value)
	} else {
		warningEmbedDescription += "but no reason was given."
	}

	if len(splittedCustomID) == 5 {
		warningEmbedDescription += fmt.Sprintf("**[This message](%s)** was attached to the warning.", fmt.Sprintf("https://discord.com/channels/%s/%s/%s", interaction.Data.GuildID, splittedCustomID[3], splittedCustomID[4]))
	}

	warningEmbed := discord.Embed{
		Description: warningEmbedDescription,
		Color:       discord.EMBED_COLOR_YELLOW,
		Image: &discord.EmbedImage{
			URL:    cdn.GetGuildIconURL(guild.ID, guild.Icon),
			Height: 100,
			Width:  100,
		},
	}

	api.SendMessage(discord.Message{
		ChannelID: dmChannel.ID,
		Embeds:    []discord.Embed{warningEmbed},
	})
}
