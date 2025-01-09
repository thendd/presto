package modals

import (
	"presto/internal/discord/api"
)

type WithHandler struct {
	Data    api.Modal
	Handler func(api.Interaction)
}

var modals = []WithHandler{}

func Get() []WithHandler {
	return modals
}

func Append(modal WithHandler) {
	modals = append(modals, modal)
}

func Remove(index int) {
	modals = append(modals[:index], modals[index+1:]...)
}
