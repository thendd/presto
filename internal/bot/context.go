package bot

import "presto/internal/discord"

type Context struct {
	Interaction discord.Interaction
	Session     *Session
}

type InteractionHandler func(context Context) error
