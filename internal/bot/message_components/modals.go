package message_components

import (
	"presto/internal/discord/api"
)

type ModalWithHandler struct {
	Data    api.Modal
	Handler func(api.Interaction)
}

var Modals = []ModalWithHandler{}
