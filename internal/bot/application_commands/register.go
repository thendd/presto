package application_commands

import (
	"slices"

	"presto/internal/discord"
	"presto/internal/discord/api"
)

type ApplicationCommandWithHandlers struct {
	Data     discord.ApplicationCommand
	Handlers []func(api.Interaction)
}

type SlashCommandGroup ApplicationCommandWithHandlers

func (group *SlashCommandGroup) AddSubCommandGroup(name, description string) *SlashCommandGroup {
	group.Data.Options = append(group.Data.Options, discord.ApplicationCommandOption{
		Type:        discord.APPLICATION_COMMAND_OPTION_TYPE_SUB_COMMAND_GROUP,
		Name:        name,
		Description: description,
	})

	return group
}

func (group *SlashCommandGroup) AddSubCommand(subCommandGroup, name, description string, options []discord.ApplicationCommandOption, handler func(api.Interaction)) *SlashCommandGroup {
	index := slices.IndexFunc(group.Data.Options, func(e discord.ApplicationCommandOption) bool {
		return e.Name == subCommandGroup && e.Type == discord.APPLICATION_COMMAND_OPTION_TYPE_SUB_COMMAND_GROUP
	})

	if index != -1 {
		group.Data.Options[index].Options = append(group.Data.Options[index].Options, discord.ApplicationCommandOption{
			Type:        discord.APPLICATION_COMMAND_OPTION_TYPE_SUB_COMMAND,
			Name:        name,
			Description: description,
			Options:     options,
		})
		group.Handlers = append(group.Handlers, handler)
	}

	return group
}

func (group *SlashCommandGroup) ToApplicationCommand() ApplicationCommandWithHandlers {
	return ApplicationCommandWithHandlers{
		Data:     group.Data,
		Handlers: group.Handlers,
	}
}

func NewSlashCommand(name, description string, options []discord.ApplicationCommandOption, handlers ...func(api.Interaction)) ApplicationCommandWithHandlers {
	return ApplicationCommandWithHandlers{
		Handlers: handlers,
		Data: discord.ApplicationCommand{
			Type:        discord.APPLICATION_COMMAND_TYPE_CHAT_INPUT,
			Name:        name,
			Description: description,
			Options:     options,
		},
	}
}

func NewSlashCommandGroup(name, description string) *SlashCommandGroup {
	return &SlashCommandGroup{
		Data: discord.ApplicationCommand{
			Name:        name,
			Description: description,
		},
	}
}

func NewUserCommand(name string, handlers ...func(api.Interaction)) ApplicationCommandWithHandlers {
	return ApplicationCommandWithHandlers{
		Handlers: handlers,
		Data: discord.ApplicationCommand{
			Type: discord.APPLICATION_COMMAND_TYPE_USER,
			Name: name,
		},
	}
}

func NewMessageCommand(name string, handlers ...func(api.Interaction)) ApplicationCommandWithHandlers {
	return ApplicationCommandWithHandlers{
		Handlers: handlers,
		Data: discord.ApplicationCommand{
			Type: discord.APPLICATION_COMMAND_TYPE_MESSAGE,
			Name: name,
		},
	}
}
