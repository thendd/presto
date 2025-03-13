package modals

import (
	"presto/internal/discord/api"
	"slices"
)

type WithHandler struct {
	Data    api.Modal
	Handler func(api.Interaction, ...any) error
	Args    []any
}

var modals = []WithHandler{}

func Get() []WithHandler {
	return modals
}

func Append(modal WithHandler) {
	modals = append(modals, modal)
}

func Remove(index int) {
	modals = slices.Delete(modals, index, index+1)
}
