package application_commands

import (
	"fmt"
	"time"

	"presto/internal/discord"
	"presto/internal/discord/api"
)

var Ping = NewSlashCommand("ping", "Have you ever heard about ping pong?", []ApplicationCommandWithHandlerDataOption{}, PingHandler).
	ToApplicationCommand()

func PingHandler(interaction api.Interaction) error {
	start := time.Now()
	// TODO: proper way of measuring the latency
	latency := time.Since(start)

	color := discord.EMBED_COLOR_RED

	if latency < 40 {
		color = discord.EMBED_COLOR_GREEN
	} else if latency <= 200 {
		color = discord.EMBED_COLOR_YELLOW
	}

	interaction.RespondWithMessage(discord.Message{
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
