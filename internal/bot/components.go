package bot

import "presto/internal/discord"

type SelectMenuWithHandler struct {
	Data    discord.MessageComponent
	Handler func(Context, ...any) error
	Args    []any
}

type ModalWithHandler struct {
	Data    discord.Modal
	Handler func(Context, ...any) error
	Args    []any
}
