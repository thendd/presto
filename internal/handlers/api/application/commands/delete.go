package commands

import (
	"log"
	"net/http"

	"presto/internal/config"
	"presto/internal/constants"
	"presto/internal/handlers/api"
)

func DeleteGlobalApplicationCommand(id string) {
	_, statusCode := api.MakeHTTPRequestToDiscord("/applications/"+config.APPLICATION_ID+"/commands/"+id, constants.METHOD_DELETE, nil)

	if statusCode != http.StatusNoContent {
		log.Fatalf("Could not delete global application command. Expected status code 204 and got %d", statusCode)
	}
}

func DeleteTestingGuildApplicationCommand(id string) {
	_, statusCode := api.MakeHTTPRequestToDiscord("/applications/"+config.APPLICATION_ID+"/guilds/"+config.TESTING_GUILD_ID+"/commands/"+id, constants.METHOD_DELETE, nil)

	if statusCode != http.StatusNoContent {
		log.Fatalf("Could not delete testing guild application command. Expected status code 204 and got %d", statusCode)
	}
}
