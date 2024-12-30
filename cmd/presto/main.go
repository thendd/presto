package main

import (
	"fmt"
	"presto/configs"
	"presto/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env")
	configs.LoadDiscordConfig()

	rawResponse := handlers.MakeHTTPRequestToDiscord("applications/"+configs.APPLICATION_ID+"/guilds/"+configs.TESTING_GUILD_ID+"/commands", configs.MethodPost, map[string]any{
		"name":        "high-five",
		"type":        1,
		"description": "Gives a high five",
	})
	fmt.Println(rawResponse)
}
