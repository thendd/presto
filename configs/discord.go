package configs

import "os"

var API_URL string
var APPLICATION_ID string
var TESTING_GUILD_ID string
var DISCORD_BOT_TOKEN string

func LoadDiscordConfig() {
	API_URL = os.Getenv("BASE_DISCORD_API_URL") + "/v" + os.Getenv("DISCORD_API_VERSION")
	APPLICATION_ID = os.Getenv("DISCORD_APPLICATION_ID")
	TESTING_GUILD_ID = os.Getenv("TESTING_GUILD_ID")
	DISCORD_BOT_TOKEN = os.Getenv("DISCORD_BOT_TOKEN")
}
