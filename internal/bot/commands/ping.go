package commands

import (
	"context"
	"fmt"
	"time"

	"presto/internal/bot/commands/contexts"
	ws "presto/internal/handlers/websocket"
	"presto/internal/types"
	"presto/tools"
)

var Ping = tools.CreateSlashCommand("ping", "Have you ever heard about ping pong?", []types.ApplicationCommandOption{}, PingHandler)

func PingHandler(ctx contexts.InteractionCreateContext) {
	start := time.Now()
	ws.Connection.Ping(context.Background())
	latency := time.Since(start).Milliseconds()

	ctx.RespondWithMessage(types.Message{
		Content: fmt.Sprintf("Pong! Latency is **%dms**", latency),
	})
}
