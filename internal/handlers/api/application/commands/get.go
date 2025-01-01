package commands

import (
	"encoding/json"
	"log"
	"net/http"

	bot_types "presto/internal/bot/types"
	"presto/internal/config"
	"presto/internal/constants"
	"presto/internal/handlers/api"
)

type getApplicationCommandsResponseBody []bot_types.ApplicationCommand

func GetGlobalApplicationCommands() (response getApplicationCommandsResponseBody) {
	rawResponse, statusCode := api.MakeHTTPRequestToDiscord("/applications/"+config.APPLICATION_ID+"/commands", constants.METHOD_GET, nil)

	if statusCode != http.StatusOK {
		log.Fatalf("Could not get global application commands. Expected status code 200 and got %d", statusCode)
	}

	json.Unmarshal(rawResponse, &response)
	return
}

func GetTestingGuildApplicationCommands() (response getApplicationCommandsResponseBody) {
	rawResponse, statusCode := api.MakeHTTPRequestToDiscord("/applications/"+config.APPLICATION_ID+"/guilds/"+config.TESTING_GUILD_ID+"/commands", constants.METHOD_GET, nil)

	if statusCode != http.StatusOK {
		log.Fatalf("Could not get testing guild application commands. Expected status code 200 and got %d", statusCode)
	}

	json.Unmarshal(rawResponse, &response)
	return
}
