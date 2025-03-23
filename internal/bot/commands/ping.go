package commands

import (
	"errors"
	"presto/internal/bot"
	"presto/internal/database"
	"presto/internal/discord"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var Ping = bot.NewSlashCommand("ping", "Have you ever heard about ping pong?", []bot.ApplicationCommandWithHandlerDataOption{}, discord.ApplicationCommandNameLocalizations{}, discord.ApplicationCommandDescriptionLocalizations{
	PtBR: "Você já ouviu falar de ping pong?",
}, PingHandler).
	ToApplicationCommand()

func PingHandler(context bot.Context) error {
	latency := context.Session.Latency.Milliseconds()

	color := discord.EMBED_COLOR_RED

	if latency < 40 {
		color = discord.EMBED_COLOR_GREEN
	} else if latency <= 200 {
		color = discord.EMBED_COLOR_YELLOW
	}

	guild := database.Guild{
		ID: context.Interaction.Data.GuildID,
	}

	database.Connection.First(&guild)

	localizer := i18n.NewLocalizer(context.Session.I18nBundle, guild.Language)
	pingText, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "Latency",
		},
		TemplateData: map[string]any{
			"Latency": latency,
		},
	})
	if err != nil {
		return errors.New("There was an error while translating the text: " + err.Error())
	}

	context.Interaction.RespondWithMessage(discord.Message{
		Embeds: []discord.Embed{
			{
				Title:       ":ping_pong: Pong!",
				Description: pingText,
				Color:       color,
			},
		},
	})

	return nil
}
