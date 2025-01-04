package application_commands

import (
	"fmt"
	"strings"
	"time"

	"presto/internal/bot/message_components"
	"presto/internal/discord"
	"presto/internal/discord/api"
	"presto/internal/discord/api/cache"
	"presto/internal/discord/cdn"
)

var Warn = NewUserCommand("Warn", WarnHandler)

func WarnHandler(interaction api.Interaction) {
	modalCustomId := fmt.Sprintf("%d-%s-%s", time.Now().UnixMilli(), interaction.Data.Member.User.ID, interaction.Data.Data.TargetID)

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
	interaction.RespondWithMessage(discord.Message{
		Embeds: []discord.Embed{
			{
				Description: "Warning was sent successfully. The user will receive a message in their DM with the reason of the warning.",
				Color:       discord.EMBED_COLOR_GREEN,
			},
		},
		Flags: discord.MESSAGE_FLAG_EPHEMERAL,
	})

	targetId := strings.Split(interaction.Data.Data.CustomID, "-")[2]
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
