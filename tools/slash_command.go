package tools

import (
	"presto/internal/bot/commands/contexts"
	bot_types "presto/internal/bot/types"
	"presto/internal/constants"
	"presto/internal/types"
)

func CreateSlashCommand(name, description string, options []types.ApplicationCommandOption, handler func(context contexts.InteractionCreateContext)) bot_types.ApplicationCommand {
	return bot_types.ApplicationCommand{
		Type:        constants.APPLICATION_COMMAND_TYPE_CHAT_INPUT,
		Name:        name,
		Description: description,
		Options:     options,
		Handler:     handler,
	}
}
