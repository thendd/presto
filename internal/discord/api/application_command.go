package api

import (
	"encoding/json"
	"net/http"
	"presto/internal/log"

	"presto/internal/discord"
	"presto/internal/discord/config"
)

const (
	APPLICATION_COMMAND_OPTION_TYPE_SUB_COMMAND discord.ApplicationCommandOptionType = iota + 1
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

type createApplicationCommandRequestBody discord.ApplicationCommand

type getApplicationCommandsResponseBody []discord.ApplicationCommand

func GetGlobalApplicationCommands() (response getApplicationCommandsResponseBody) {
	rawResponse, statusCode := MakeRequest("/applications/"+config.APPLICATION_ID+"/commands", http.MethodGet, nil)

	if statusCode != http.StatusOK {
		log.Fatal("Could not get global application commands. Expected status code 200 and got %d", statusCode)
	}

	json.Unmarshal(rawResponse, &response)
	return
}

func GetTestingGuildApplicationCommands() (response getApplicationCommandsResponseBody) {
	rawResponse, statusCode := MakeRequest("/applications/"+config.APPLICATION_ID+"/guilds/"+config.TESTING_GUILD_ID+"/commands", http.MethodGet, nil)

	if statusCode != http.StatusOK {
		log.Fatal("Could not get testing guild application commands. Expected status code 200 and got %d", statusCode)
	}

	json.Unmarshal(rawResponse, &response)
	return
}

func CreateGlobalApplicationCommand(applicationCommand discord.ApplicationCommand) {
	body, _ := json.Marshal(applicationCommand)
	_, statusCode := MakeRequest("/applications/"+config.APPLICATION_ID+"/commands", http.MethodPost, body)

	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		log.Fatal("Could not create global application command. Expected status code 200 or 201 and got %d", statusCode)
	}
}

func CreateTestingGuildApplicationCommand(applicationCommand discord.ApplicationCommand) {
	_, statusCode := MakeRequest("/applications/"+config.APPLICATION_ID+"/guilds/"+config.TESTING_GUILD_ID+"/commands", http.MethodPost, applicationCommand)

	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		log.Fatal("Could not create testing guild application commands. Expected status code 200 or 201 and got %d", statusCode)
	}
}

func DeleteGlobalApplicationCommand(id string) {
	_, statusCode := MakeRequest("/applications/"+config.APPLICATION_ID+"/commands/"+id, http.MethodDelete, nil)

	if statusCode != http.StatusNoContent {
		log.Fatal("Could not delete global application command. Expected status code 204 and got %d", statusCode)
	}
}

func DeleteTestingGuildApplicationCommand(id string) {
	_, statusCode := MakeRequest("/applications/"+config.APPLICATION_ID+"/guilds/"+config.TESTING_GUILD_ID+"/commands/"+id, http.MethodDelete, nil)

	if statusCode != http.StatusNoContent {
		log.Fatal("Could not delete testing guild application command. Expected status code 204 and got %d", statusCode)
	}
}
