package application_commands

import (
	"context"
	"fmt"
	"time"

	"presto/internal/discord"
	"presto/internal/discord/api"
)

var Ping = NewSlashCommand("ping", "Have you ever heard about ping pong?", []discord.ApplicationCommandOption{}, PingHandler)

func PingHandler(interaction api.Interaction) error {
	start := time.Now()
	interaction.Websocket.Ping(context.Background())
	latency := time.Since(start).Milliseconds()

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
