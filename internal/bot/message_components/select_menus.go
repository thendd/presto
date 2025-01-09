package message_components

import (
	"presto/internal/discord"
	"presto/internal/discord/api"
)

type SelectMenuWithHandler struct {
	Data                  discord.MessageComponent
	DeleteAfterInteracted bool
	Handler               func(api.Interaction)
}

var SelectMenus = []SelectMenuWithHandler{}
