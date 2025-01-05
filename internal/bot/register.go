package bot

import (
	"log"
	"os"
	"slices"

	"presto/internal/bot/application_commands"
	"presto/internal/discord"
	"presto/internal/discord/api"
)

var RegisteredCommands = []application_commands.ApplicationCommandWithHandler{application_commands.Ping, application_commands.WarnUserCommand, application_commands.WarnSlashCommand}

func RegisterCommands() {
	log.Println("Started registering application commands")

	var mustDelete []discord.ApplicationCommand
	var mustCreate []discord.ApplicationCommand

	switch os.Getenv("ENVIRONMENT") {
	case "production":
		applicationCommands := api.GetGlobalApplicationCommands()

		for _, registeredCommand := range RegisteredCommands {
			possibleExistingApplicationCommandIndex := slices.IndexFunc(applicationCommands, func(e discord.ApplicationCommand) bool {
				return e.Name == registeredCommand.Data.Name
			})

			if possibleExistingApplicationCommandIndex != -1 {
				areRegisteredCommandAndApplicationCommandEqual := discord.CompareApplicationCommands(registeredCommand.Data, applicationCommands[possibleExistingApplicationCommandIndex])

				if !areRegisteredCommandAndApplicationCommandEqual {
					// Created application commands with names that are already assigned to some application command are overwritten
					mustCreate = append(mustCreate, registeredCommand.Data)
				}

				applicationCommands = append(applicationCommands[:possibleExistingApplicationCommandIndex], applicationCommands[possibleExistingApplicationCommandIndex+1:]...)
			} else {
				mustCreate = append(mustCreate, registeredCommand.Data)
			}
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
				return e.Name == registeredCommand.Data.Name
			})

			if possibleExistingApplicationCommandIndex != -1 {
				areRegisteredCommandAndApplicationCommandEqual := discord.CompareApplicationCommands(registeredCommand.Data, applicationCommands[possibleExistingApplicationCommandIndex])

				if !areRegisteredCommandAndApplicationCommandEqual {
					// Created application commands with names that are already assigned to some application command are overwritten
					mustCreate = append(mustCreate, registeredCommand.Data)
				}

				applicationCommands = append(applicationCommands[:possibleExistingApplicationCommandIndex], applicationCommands[possibleExistingApplicationCommandIndex+1:]...)
			} else {
				mustCreate = append(mustCreate, registeredCommand.Data)
			}
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
		log.Fatalf("Unknown \"ENVIRONMENT\" value: %v", os.Getenv("ENVIRONMENT"))
	}

	log.Println("Finished registering commands successfully")
}
