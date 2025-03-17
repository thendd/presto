package bot

import (
	"os"
	"presto/internal/discord"
	"presto/internal/log"
	"slices"
)

type ApplicationCommandWithHandlerDataOptionChoice struct {
	Name          string            `json:"name"`
	Localizations map[string]string `json:"name_localizations"`
	Value         string            `json:"value"`
}

type ApplicationCommandWithHandlerDataOption struct {
	Type         discord.ApplicationCommandOptionType            `json:"type"`
	Name         string                                          `json:"name"`
	Description  string                                          `json:"description"`
	Required     bool                                            `json:"required"`
	Autocomplete bool                                            `json:"autocomplete"`
	Options      []ApplicationCommandWithHandlerDataOption       `json:"options,omitempty"`
	Choices      []ApplicationCommandWithHandlerDataOptionChoice `json:"choices,omitempty"`
	Handler      InteractionHandler                              `json:"-"`
}

type ApplicationCommandWithHandlerData struct {
	ID          any                                       `json:"id,omitempty"`
	Name        string                                    `json:"name"`
	Description string                                    `json:"description"`
	Options     []ApplicationCommandWithHandlerDataOption `json:"options,omitempty"`
	Type        discord.ApplicationCommandType            `json:"type,omitempty"`
}

type ApplicationCommandWithHandler struct {
	Data    ApplicationCommandWithHandlerData
	Handler InteractionHandler
}

func (applicationCommandWithHandler ApplicationCommandWithHandler) ToApplicationCommand() discord.ApplicationCommand {
	var lookForOptions func(options []ApplicationCommandWithHandlerDataOption) []discord.ApplicationCommandOption
	lookForOptions = func(options []ApplicationCommandWithHandlerDataOption) []discord.ApplicationCommandOption {
		var final []discord.ApplicationCommandOption

		for _, option := range options {
			finalOptions := lookForOptions(option.Options)

			final = append(final, discord.ApplicationCommandOption{
				Type:         option.Type,
				Name:         option.Name,
				Description:  option.Description,
				Required:     option.Required,
				Autocomplete: option.Autocomplete,
				Options:      finalOptions,
			})
		}

		return final
	}

	return discord.ApplicationCommand{
		ID:          applicationCommandWithHandler.Data.ID,
		Name:        applicationCommandWithHandler.Data.Name,
		Description: applicationCommandWithHandler.Data.Description,
		Type:        applicationCommandWithHandler.Data.Type,
		Options:     lookForOptions(applicationCommandWithHandler.Data.Options),
	}
}

type SlashCommand ApplicationCommandWithHandler

func (command *SlashCommand) AddSubCommandGroup(name string) *SlashCommand {
	command.Data.Options = append(command.Data.Options, ApplicationCommandWithHandlerDataOption{
		Type: discord.APPLICATION_COMMAND_OPTION_TYPE_SUB_COMMAND_GROUP,
		Name: name,
	})

	return command
}

func (command *SlashCommand) AddSubCommand(name, description string, options []ApplicationCommandWithHandlerDataOption, handler InteractionHandler) *SlashCommand {
	command.Data.Options = append(command.Data.Options, ApplicationCommandWithHandlerDataOption{
		Type:        discord.APPLICATION_COMMAND_OPTION_TYPE_SUB_COMMAND,
		Name:        name,
		Description: description,
		Options:     options,
		Handler:     handler,
	})

	return command
}

func (command *SlashCommand) ToApplicationCommand() ApplicationCommandWithHandler {
	return ApplicationCommandWithHandler{
		Data:    command.Data,
		Handler: command.Handler,
	}
}

func NewSlashCommand(name, description string, options []ApplicationCommandWithHandlerDataOption, handler InteractionHandler) *SlashCommand {
	return &SlashCommand{
		Handler: handler,
		Data: ApplicationCommandWithHandlerData{
			Type:        discord.APPLICATION_COMMAND_TYPE_CHAT_INPUT,
			Name:        name,
			Description: description,
			Options:     options,
		},
	}
}

func NewUserCommand(name string, handler InteractionHandler) ApplicationCommandWithHandler {
	return ApplicationCommandWithHandler{
		Handler: handler,
		Data: ApplicationCommandWithHandlerData{
			Type: discord.APPLICATION_COMMAND_TYPE_USER,
			Name: name,
		},
	}
}

func NewMessageCommand(name string, handler InteractionHandler) ApplicationCommandWithHandler {
	return ApplicationCommandWithHandler{
		Handler: handler,
		Data: ApplicationCommandWithHandlerData{
			Type: discord.APPLICATION_COMMAND_TYPE_MESSAGE,
			Name: name,
		},
	}
}

func PushCommands(commands []ApplicationCommandWithHandler) {
	log.Info("Started registering application commands")

	mustDelete := []discord.ApplicationCommand{}
	mustCreate := []discord.ApplicationCommand{}

	var localApplicationCommands []discord.ApplicationCommand

	for _, localApplicationCommand := range commands {
		localApplicationCommands = append(localApplicationCommands, localApplicationCommand.ToApplicationCommand())
	}

	var applicationCommands []discord.ApplicationCommand

	switch os.Getenv("PRESTO_ENVIRONMENT") {
	case "production":
		applicationCommands = discord.GetGlobalApplicationCommands()
	case "development":
		applicationCommands = discord.GetTestingGuildApplicationCommands()
	default:
		log.Fatalf("Unknown \"PRESTO_ENVIRONMENT\" value: %v", os.Getenv("PRESTO_ENVIRONMENT"))
	}

	for _, applicationCommand := range applicationCommands {
		exactMatchIndex := slices.IndexFunc(localApplicationCommands, func(localApplicationCommand discord.ApplicationCommand) bool {
			return discord.CompareApplicationCommands(applicationCommand, localApplicationCommand)
		})

		if exactMatchIndex == -1 {
			mustDelete = append(mustDelete, applicationCommand)
			continue
		}

		localApplicationCommands = slices.Delete(localApplicationCommands, exactMatchIndex, exactMatchIndex+1)
	}

	mustCreate = localApplicationCommands

	switch os.Getenv("PRESTO_ENVIRONMENT") {
	case "production":
		for _, applicationCommand := range mustDelete {
			discord.DeleteGlobalApplicationCommand(applicationCommand.ID.(string))
			log.Infof("\"%s\" command was deleted globally and successfully", applicationCommand.Name)
		}

		for _, applicationCommand := range mustCreate {
			discord.CreateGlobalApplicationCommand(applicationCommand)
			log.Infof("\"%s\" command was created/updated globaly successfully", applicationCommand.Name)
		}
	case "development":
		for _, applicationCommand := range mustDelete {
			discord.DeleteTestingGuildApplicationCommand(applicationCommand.ID.(string))
			log.Infof("\"%s\" command was deleted successfully in the testing guild", applicationCommand.Name)
		}

		for _, applicationCommand := range mustCreate {
			discord.CreateTestingGuildApplicationCommand(applicationCommand)
			log.Infof("\"%s\" command was created/updated successfully in the testing guild", applicationCommand.Name)
		}
	}

	log.Info("Finished registering commands successfully")
}
