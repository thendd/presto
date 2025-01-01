package types

import (
	"presto/internal/bot/commands/contexts"
	"presto/internal/types"
)

type ApplicationCommand struct {
	ID          any                                             `json:"id,omitempty"`
	Name        string                                          `json:"name"`
	Description string                                          `json:"description"`
	Options     []types.ApplicationCommandOption                `json:"options,omitempty"`
	Type        types.ApplicationCommandType                    `json:"type"`
	Handler     func(context contexts.InteractionCreateContext) `json:"-"`
}
