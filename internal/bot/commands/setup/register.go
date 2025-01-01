package setup

import (
	"log"
	"os"
	"slices"

	"presto/internal/bot/commands"
	bot_types "presto/internal/bot/types"
	application_commands "presto/internal/handlers/api/application/commands"
	"presto/tools"
)

var RegisteredCommands = []bot_types.ApplicationCommand{commands.Ping}

func RegisterCommands() {
	var mustDelete []bot_types.ApplicationCommand
	var mustCreate []bot_types.ApplicationCommand

	switch os.Getenv("ENVIRONMENT") {
	case "production":
		applicationCommands := application_commands.GetGlobalApplicationCommands()

		for _, registeredCommand := range RegisteredCommands {
			possibleExistingApplicationCommandIndex := slices.IndexFunc(applicationCommands, func(e bot_types.ApplicationCommand) bool {
				return e.Name == registeredCommand.Name
			})

			if possibleExistingApplicationCommandIndex != -1 {
				areRegisteredCommandAndApplicationCommandEqual := tools.CompareApplicationCommands(registeredCommand, applicationCommands[possibleExistingApplicationCommandIndex])

				if !areRegisteredCommandAndApplicationCommandEqual {
					// Created application commands with names that are already assigned to some application command are overwritten
					mustCreate = append(mustCreate, registeredCommand)
				}

				applicationCommands = append(applicationCommands[:possibleExistingApplicationCommandIndex], applicationCommands[possibleExistingApplicationCommandIndex+1:]...)
			} else {
				mustCreate = append(mustCreate, registeredCommand)
			}
		}

		// If an application command exists and does not correspond with the name of any of the registered commands,
		// it must be deleted
		mustDelete = applicationCommands

		for _, applicationCommand := range mustDelete {
			application_commands.DeleteGlobalApplicationCommand(applicationCommand.ID.(string))
			log.Printf("\"%s\" command was deleted globally and successfully\n", applicationCommand.Name)
		}

		for _, applicationCommand := range mustCreate {
			application_commands.CreateGlobalApplicationCommand(applicationCommand)
			log.Printf("\"%s\" command was created/updated globaly successfully\n", applicationCommand.Name)
		}
	case "development":
		applicationCommands := application_commands.GetTestingGuildApplicationCommands()

		for _, registeredCommand := range RegisteredCommands {
			possibleExistingApplicationCommandIndex := slices.IndexFunc(applicationCommands, func(e bot_types.ApplicationCommand) bool {
				return e.Name == registeredCommand.Name
			})

			if possibleExistingApplicationCommandIndex != -1 {
				areRegisteredCommandAndApplicationCommandEqual := tools.CompareApplicationCommands(registeredCommand, applicationCommands[possibleExistingApplicationCommandIndex])

				if !areRegisteredCommandAndApplicationCommandEqual {
					// Created application commands with names that are already assigned to some application command are overwritten
					mustCreate = append(mustCreate, registeredCommand)
				}

				applicationCommands = append(applicationCommands[:possibleExistingApplicationCommandIndex], applicationCommands[possibleExistingApplicationCommandIndex+1:]...)
			} else {
				mustCreate = append(mustCreate, registeredCommand)
			}
		}

		// If an application command exists and does not correspond with the name of any of the registered commands,
		// it must be deleted
		mustDelete = applicationCommands

		for _, applicationCommand := range mustDelete {
			application_commands.DeleteTestingGuildApplicationCommand(applicationCommand.ID.(string))
			log.Printf("\"%s\" command was deleted successfully in the testing guild\n", applicationCommand.Name)
		}

		for _, applicationCommand := range mustCreate {
			application_commands.CreateTestingGuildApplicationCommand(applicationCommand)
			log.Printf("\"%s\" command was created/updated successfully in the testing guild\n", applicationCommand.Name)
		}
	default:
		log.Fatalf("Unknown \"ENVIRONMENT\" value: %v", os.Getenv("ENVIRONMENT"))
	}
}
