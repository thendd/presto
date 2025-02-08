package application_commands

import (
	"log"
	"os"
	"slices"

	"presto/internal/discord"
	"presto/internal/discord/api"
)

type ApplicationCommandWithHandlers struct {
	Data     discord.ApplicationCommand
	Handlers []func(api.Interaction) error
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

func (group *SlashCommandGroup) AddSubCommand(subCommandGroup, name, description string, options []discord.ApplicationCommandOption, handler func(api.Interaction) error) *SlashCommandGroup {
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

func NewSlashCommand(name, description string, options []discord.ApplicationCommandOption, handlers ...func(api.Interaction) error) ApplicationCommandWithHandlers {
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

func NewUserCommand(name string, handlers ...func(api.Interaction) error) ApplicationCommandWithHandlers {
	return ApplicationCommandWithHandlers{
		Handlers: handlers,
		Data: discord.ApplicationCommand{
			Type: discord.APPLICATION_COMMAND_TYPE_USER,
			Name: name,
		},
	}
}

func NewMessageCommand(name string, handlers ...func(api.Interaction) error) ApplicationCommandWithHandlers {
	return ApplicationCommandWithHandlers{
		Handlers: handlers,
		Data: discord.ApplicationCommand{
			Type: discord.APPLICATION_COMMAND_TYPE_MESSAGE,
			Name: name,
		},
	}
}

var RegisteredCommands = []ApplicationCommandWithHandlers{
	Ping,
	WarnUserCommand,
	WarnSlashCommand,
	WarnMessageCommand,
	Settings,
}

func Register() {
	log.Println("Started registering application commands")

	mustDelete := []discord.ApplicationCommand{}
	mustCreate := []discord.ApplicationCommand{}

	switch os.Getenv("PRESTO_ENVIRONMENT") {
	case "production":
		applicationCommands := api.GetGlobalApplicationCommands()

		for _, registeredCommand := range RegisteredCommands {
			possibleExistingApplicationCommandIndex := slices.IndexFunc(applicationCommands, func(e discord.ApplicationCommand) bool {
				return discord.GetApplicationCommandName(registeredCommand.Data) == discord.GetApplicationCommandName(e)
			})

			if possibleExistingApplicationCommandIndex == -1 || !discord.CompareApplicationCommands(registeredCommand.Data, applicationCommands[possibleExistingApplicationCommandIndex]) {
				mustCreate = append(mustCreate, registeredCommand.Data)
				return
			}

			applicationCommands = append(applicationCommands[:possibleExistingApplicationCommandIndex], applicationCommands[possibleExistingApplicationCommandIndex+1:]...)
		}

		// If an application command exists and does not correspond with the name of any of the registered commands,
		// it must be deleted
		mustDelete = applicationCommands

		for _, applicationCommand := range mustDelete {
			api.DeleteGlobalApplicationCommand(applicationCommand.ID.(string))
			log.Printf("\"%s\" command was deleted globally and successfully\n", applicationCommand.Name)
		}

		for _, applicationCommand := range mustCreate {
			api.CreateGlobalApplicationCommand(applicationCommand)
			log.Printf("\"%s\" command was created/updated globaly successfully\n", applicationCommand.Name)
		}
	case "development":
		applicationCommands := api.GetTestingGuildApplicationCommands()

		for _, registeredCommand := range RegisteredCommands {
			possibleExistingApplicationCommandIndex := slices.IndexFunc(applicationCommands, func(e discord.ApplicationCommand) bool {
				return discord.GetApplicationCommandName(registeredCommand.Data) == discord.GetApplicationCommandName(e)
			})

			if possibleExistingApplicationCommandIndex == -1 || !discord.CompareApplicationCommands(registeredCommand.Data, applicationCommands[possibleExistingApplicationCommandIndex]) {
				mustCreate = append(mustCreate, registeredCommand.Data)
				return
			}

			applicationCommands = append(applicationCommands[:possibleExistingApplicationCommandIndex], applicationCommands[possibleExistingApplicationCommandIndex+1:]...)
		}

		// If an application command exists and does not correspond with the name of any of the registered commands,
		// it must be deleted
		mustDelete = applicationCommands

		for _, applicationCommand := range mustDelete {
			api.DeleteTestingGuildApplicationCommand(applicationCommand.ID.(string))
			log.Printf("\"%s\" command was deleted successfully in the testing guild\n", applicationCommand.Name)
		}

		for _, applicationCommand := range mustCreate {
			api.CreateTestingGuildApplicationCommand(applicationCommand)
			log.Printf("\"%s\" command was created/updated successfully in the testing guild\n", applicationCommand.Name)
		}
	default:
		log.Fatalf("Unknown \"PRESTO_ENVIRONMENT\" value: %v", os.Getenv("PRESTO_ENVIRONMENT"))
	}

	log.Println("Finished registering commands successfully")
}
