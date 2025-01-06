package application_commands

import (
	"presto/internal/discord"
	"presto/internal/discord/api"
)

type ApplicationCommandWithHandler struct {
	Data    discord.ApplicationCommand
	Handler func(api.Interaction)
}

var SlashCommands = []ApplicationCommandWithHandler{}

func NewSlashCommand(name, description string, options []discord.ApplicationCommandOption, handler func(api.Interaction)) ApplicationCommandWithHandler {
	return ApplicationCommandWithHandler{
		Handler: handler,
		Data: discord.ApplicationCommand{
			Type:        discord.APPLICATION_COMMAND_TYPE_CHAT_INPUT,
			Name:        name,
			Description: description,
			Options:     options,
		},
	}
}

func NewUserCommand(name string, handler func(api.Interaction)) ApplicationCommandWithHandler {
	return ApplicationCommandWithHandler{
		Handler: handler,
		Data: discord.ApplicationCommand{
			Type: discord.APPLICATION_COMMAND_TYPE_USER,
			Name: name,
		},
	}
}

func NewMessageCommand(name string, handler func(api.Interaction)) ApplicationCommandWithHandler {
	return ApplicationCommandWithHandler{
		Handler: handler,
		Data: discord.ApplicationCommand{
			Type: discord.APPLICATION_COMMAND_TYPE_MESSAGE,
			Name: name,
		},
	}
}
