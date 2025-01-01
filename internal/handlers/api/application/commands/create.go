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

type CreateApplicationCommandRequestBody bot_types.ApplicationCommand

func CreateGlobalApplicationCommand(applicationCommand bot_types.ApplicationCommand) {
	body, _ := json.Marshal(applicationCommand)
	_, statusCode := api.MakeHTTPRequestToDiscord("/applications/"+config.APPLICATION_ID+"/commands", constants.METHOD_POST, body)

	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		log.Fatalf("Could not create global application command. Expected status code 200 or 201 and got %d", statusCode)
	}
}

func CreateTestingGuildApplicationCommand(applicationCommand bot_types.ApplicationCommand) {
	_, statusCode := api.MakeHTTPRequestToDiscord("/applications/"+config.APPLICATION_ID+"/guilds/"+config.TESTING_GUILD_ID+"/commands", constants.METHOD_POST, applicationCommand)

	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		log.Fatalf("Could not create testing guild application commands. Expected status code 200 or 201 and got %d", statusCode)
	}
}
