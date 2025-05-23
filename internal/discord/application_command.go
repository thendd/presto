package discord

import (
	"encoding/json"
	"net/http"
	"presto/internal/config"
	"presto/internal/log"
	"slices"
)

const (
	APPLICATION_COMMAND_TYPE_CHAT_INPUT ApplicationCommandType = iota + 1
	APPLICATION_COMMAND_TYPE_USER
	APPLICATION_COMMAND_TYPE_MESSAGE
	APPLICATION_COMMAND_TYPE_PRIMARY_ENTRY_POINT
)

const (
	APPLICATION_COMMAND_OPTION_TYPE_SUB_COMMAND ApplicationCommandOptionType = iota + 1
	APPLICATION_COMMAND_OPTION_TYPE_SUB_COMMAND_GROUP
	APPLICATION_COMMAND_OPTION_TYPE_STRING
	APPLICATION_COMMAND_OPTION_TYPE_INTEGER
	APPLICATION_COMMAND_OPTION_TYPE_BOOLEAN
	APPLICATION_COMMAND_OPTION_TYPE_USER
	APPLICATION_COMMAND_OPTION_TYPE_CHANNEL
	APPLICATION_COMMAND_OPTION_TYPE_ROLE
	APPLICATION_COMMAND_OPTION_TYPE_MENTIONABLE
	APPLICATION_COMMAND_OPTION_TYPE_NUMBER
	APPLICATION_COMMAND_OPTION_TYPE_ATTACHMENT
)

type (
	ApplicationCommandType       int
	ApplicationCommandOptionType int
)

type ApplicationCommand struct {
	ID          any                        `json:"id,omitempty"`
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Options     []ApplicationCommandOption `json:"options,omitempty"`
	Type        ApplicationCommandType     `json:"type,omitempty"`
}

type ApplicationCommandOption struct {
	Type         ApplicationCommandOptionType `json:"type"`
	Name         string                       `json:"name"`
	Description  string                       `json:"description"`
	Required     bool                         `json:"required"`
	Autocomplete bool                         `json:"autocomplete"`
	Options      []ApplicationCommandOption   `json:"options,omitempty"`
	MinimumValue int                          `json:"min_value,omitempty"`
	MaximumValue int                          `json:"max_value,omitempty"`
}

func GetFullNamesOfApplicationCommand(command ApplicationCommand) []string {
	var names []string

	for _, option := range command.Options {
		if option.Type == APPLICATION_COMMAND_OPTION_TYPE_SUB_COMMAND {
			names = append(names, command.Name+" "+option.Name)
		}
	}

	if len(names) == 0 {
		names = append(names, command.Name)
	}

	return names
}

// Compares two `[]ApplicationCommandOption` and returns whether they are exaclty equal or not
func AreApplicationCommandOptionsEqual(a []ApplicationCommandOption, b []ApplicationCommandOption) bool {
	if len(a) != len(b) {
		return false
	}

	for i, x := range a {
		y := a[i]

		if x.Type != x.Type || x.Description != y.Description || x.Required != y.Required || x.Autocomplete != y.Autocomplete {
			return false
		}

		if !AreApplicationCommandOptionsEqual(x.Options, y.Options) {
			return false
		}
	}

	return true
}

// Compares two application commands without taking the ID into consideration.
// This is because registered commands are structs created with the necessary data
// to create or update an application command and the ID is generated by Discord so
// every comparision using `reflect.DeepEqual` would fail since the command registered locally
// will have `nil` as the ID while the application command fetched from Discord will have some
// snowflake as the ID.
func CompareApplicationCommands(a ApplicationCommand, b ApplicationCommand) bool {
	if !slices.Equal(GetFullNamesOfApplicationCommand(a), GetFullNamesOfApplicationCommand(b)) {
		return false
	}

	return a.Description == b.Description && a.Type == b.Type && AreApplicationCommandOptionsEqual(a.Options, b.Options)
}

// Fetches global application commands
func GetGlobalApplicationCommands() (response []ApplicationCommand) {
	rawResponse, statusCode := MakeRequest("/applications/"+config.DISCORD_APPLICATION_ID+"/commands", http.MethodGet, nil, map[string]string{})

	if statusCode != http.StatusOK {
		log.Fatal("Could not get global application commands. Expected status code 200 and got %d", statusCode)
	}

	json.Unmarshal(rawResponse, &response)
	return
}

// Fetches application commands registered in the test guild
func GetTestingGuildApplicationCommands() (response []ApplicationCommand) {
	rawResponse, statusCode := MakeRequest("/applications/"+config.DISCORD_APPLICATION_ID+"/guilds/"+config.DISCORD_TESTING_GUILD_ID+"/commands", http.MethodGet, nil, map[string]string{})

	if statusCode != http.StatusOK {
		log.Fatal("Could not get testing guild application commands. Expected status code 200 and got %d", statusCode)
	}

	json.Unmarshal(rawResponse, &response)
	return
}

// Creates a global application command
func CreateGlobalApplicationCommand(applicationCommand ApplicationCommand) {
	body, _ := json.Marshal(applicationCommand)
	_, statusCode := MakeRequest("/applications/"+config.DISCORD_APPLICATION_ID+"/commands", http.MethodPost, body, map[string]string{})

	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		log.Fatal("Could not create global application command. Expected status code 200 or 201 and got %d", statusCode)
	}
}

// Creates an application command in the testing guild
func CreateTestingGuildApplicationCommand(applicationCommand ApplicationCommand) {
	_, statusCode := MakeRequest("/applications/"+config.DISCORD_APPLICATION_ID+"/guilds/"+config.DISCORD_TESTING_GUILD_ID+"/commands", http.MethodPost, applicationCommand, map[string]string{})

	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		log.Fatal("Could not create testing guild application commands. Expected status code 200 or 201 and got %d", statusCode)
	}
}

// Deletes an application command globally
func DeleteGlobalApplicationCommand(id string) {
	_, statusCode := MakeRequest("/applications/"+config.DISCORD_APPLICATION_ID+"/commands/"+id, http.MethodDelete, nil, map[string]string{})

	if statusCode != http.StatusNoContent {
		log.Fatal("Could not delete global application command. Expected status code 204 and got %d", statusCode)
	}
}

// Deletes an application command in the testing guild
func DeleteTestingGuildApplicationCommand(id string) {
	_, statusCode := MakeRequest("/applications/"+config.DISCORD_APPLICATION_ID+"/guilds/"+config.DISCORD_TESTING_GUILD_ID+"/commands/"+id, http.MethodDelete, nil, map[string]string{})

	if statusCode != http.StatusNoContent {
		log.Fatal("Could not delete testing guild application command. Expected status code 204 and got %d", statusCode)
	}
}
