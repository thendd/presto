package message_components

import (
	"presto/internal/discord"
	"presto/internal/discord/api"
)

type SelectMenuWithHandler struct {
	Data    discord.MessageComponent
	Handler func(api.Interaction, ...any) error
	Args    []any
}

var SelectMenus = []SelectMenuWithHandler{}
