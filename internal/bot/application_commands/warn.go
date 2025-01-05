package application_commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"presto/internal/bot/message_components"
	"presto/internal/database"
	"presto/internal/discord"
	"presto/internal/discord/api"
	"presto/internal/discord/api/cache"
	"presto/internal/discord/cdn"

	"github.com/jackc/pgx/v5/pgtype"
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
)

func WarnHandler(interaction api.Interaction) {
	modalCustomId := fmt.Sprintf("%d-%s-", time.Now().UnixMilli(), interaction.Data.Member.User.ID)

	switch interaction.Data.Data.Type {
	case discord.APPLICATION_COMMAND_TYPE_CHAT_INPUT:
		modalCustomId += interaction.Data.Data.Options[0].Value.(string)
	case discord.APPLICATION_COMMAND_TYPE_USER:
		modalCustomId += interaction.Data.Data.TargetID
	}

	modal := message_components.ModalWithHandler{
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
							Label:     "Reason for the warning",
							Style:     discord.TEXT_INPUT_STYLE_PARAGRAPH,
							MaxLength: 1500,
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
	targetId := strings.Split(interaction.Data.Data.CustomID, "-")[2]

	queries := database.New(database.Connection)
	guildData, _ := queries.GetGuild(context.Background(), interaction.Data.GuildID)
	userWarnings, err := queries.GetWarningsFromWarnedUser(context.Background(), database.GetWarningsFromWarnedUserParams{GuildID: interaction.Data.GuildID, UserID: targetId})
	if err != nil {
		userWarnings, _ = queries.CreateWarnedUser(context.Background(), database.CreateWarnedUserParams{
			GuildID: interaction.Data.GuildID,
			UserID:  targetId,
		})
	}

	remainingWarnings := guildData.MaxWarningsPerUser.Int32 - userWarnings.Int32 - 1

	interaction.RespondWithMessage(discord.Message{
		Embeds: []discord.Embed{
			{
				Description: fmt.Sprintf("Warning was sent successfully. The user will receive a message in their DM with the reason of the warning. They still have **%d warnings left**.", remainingWarnings),
				Color:       discord.EMBED_COLOR_GREEN,
			},
		},
		Flags: discord.MESSAGE_FLAG_EPHEMERAL,
	})

	// Not sure whether this is the best way to do this, but if it's working, it's working
	x := pgtype.Int4{
		Int32: userWarnings.Int32 + 1,
		Valid: true,
	}

	queries.UpdateWarnedUserWarnings(context.Background(), database.UpdateWarnedUserWarningsParams{
		Warnings: x,
		GuildID:  interaction.Data.GuildID,
		UserID:   targetId,
	})

	dmChannel := cache.GetDMChannelByRecipientID(targetId)
	if dmChannel.ID == "" {
		dmChannel = api.CreateDM(targetId)
	}

	guild := cache.GetGuildById(interaction.Data.GuildID)
	if guild.Name == "" {
		guild = api.GetGuildById(interaction.Data.GuildID)
	}

	api.SendMessage(discord.Message{
		ChannelID: dmChannel.ID,
		Embeds: []discord.Embed{
			{
				Title: "You received an warning!",
				Color: discord.EMBED_COLOR_YELLOW,
				Image: &discord.EmbedImage{
					URL:    cdn.GetGuildIconURL(guild.ID, guild.Icon),
					Height: 100,
					Width:  100,
				},
				Fields: []discord.EmbedField{
					{
						Name:  "Server name",
						Value: guild.Name,
					},
					{
						Name:  "Reason",
						Value: interaction.Data.Data.Components[0].Components[0].Value,
					},
				},
			},
		},
	})
}
