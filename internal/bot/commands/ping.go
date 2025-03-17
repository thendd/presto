package commands

import (
	"fmt"
	"presto/internal/bot"
	"presto/internal/discord"
)

var Ping = bot.NewSlashCommand("ping", "Have you ever heard about ping pong?", []bot.ApplicationCommandWithHandlerDataOption{}, PingHandler).
	ToApplicationCommand()

func PingHandler(context bot.Context) error {
	latency := context.Session.Latency.Milliseconds()

	color := discord.EMBED_COLOR_RED

	if latency < 40 {
		color = discord.EMBED_COLOR_GREEN
	} else if latency <= 200 {
		color = discord.EMBED_COLOR_YELLOW
	}

	context.Interaction.RespondWithMessage(discord.Message{
		Embeds: []discord.Embed{
			{
				Title:       ":ping_pong: Pong!",
				Description: fmt.Sprintf("Latency is **%dms**", latency),
				Color:       color,
			},
		},
	})

	return nil
}
