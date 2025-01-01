package commands

import (
	"presto/internal/bot/commands/contexts"
	"presto/internal/types"
	"presto/tools"
)

var Ping = tools.CreateSlashCommand("ping", "Have you ever heard about ping pong?", []types.ApplicationCommandOption{}, PingHandler)

func PingHandler(context contexts.InteractionCreateContext) {
	context.RespondWithMessage(types.Message{
		Content: "Pong!",
	})
}
