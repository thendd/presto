package commands

import (
	"context"
	"fmt"
	"time"

	"presto/internal/bot/commands/contexts"
	"presto/internal/constants"
	ws "presto/internal/handlers/websocket"
	"presto/internal/types"
	"presto/tools"
)

var Ping = tools.CreateSlashCommand("ping", "Have you ever heard about ping pong?", []types.ApplicationCommandOption{}, PingHandler)

func PingHandler(ctx contexts.InteractionCreateContext) {
	start := time.Now()
	ws.Connection.Ping(context.Background())
	latency := time.Since(start).Milliseconds()

	color := constants.EMBED_COLOR_RED

	if latency < 40 {
		color = constants.EMBED_COLOR_GREEN
	} else if latency <= 200 {
		color = constants.EMBED_COLOR_YELLOW
	}

	ctx.RespondWithMessage(types.Message{
		Embeds: []types.Embed{
			{
				Title:       ":ping_pong: Pong!",
				Description: fmt.Sprintf("Latency is **%dms**", latency),
				Color:       color,
			},
		},
	})
}
