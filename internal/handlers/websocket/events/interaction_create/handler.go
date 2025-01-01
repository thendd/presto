package interaction_create

import (
	"encoding/json"
	"slices"

	"presto/internal/bot/commands/contexts"
	"presto/internal/bot/commands/setup"
	bot_types "presto/internal/bot/types"
)

func HandleInteractionCreate(data []byte) {
	var context contexts.InteractionCreateContext
	json.Unmarshal(data, &context)

	index := slices.IndexFunc(setup.RegisteredCommands, func(e bot_types.ApplicationCommand) bool {
		return e.Name == context.Data.Name
	})

	setup.RegisteredCommands[index].Handler(context)
}
