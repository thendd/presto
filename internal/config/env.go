package config

import (
	"os"
	"presto/internal/log"

	"github.com/joho/godotenv"
)

var (
	POSTGRESQL_CONNECTION_STRING string
	DISCORD_APPLICATION_ID       string
	DISCORD_TESTING_GUILD_ID     string
	DISCORD_BOT_TOKEN            string
	DISCORD_API_BASE_URL         string
	DISCORD_CDN_BASE_URL         string
)

func LoadEnvironmentVariables() {
	log.Info("Started loading environment variables")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	DISCORD_API_BASE_URL = os.Getenv("BASE_DISCORD_API_URL") + "/v" + os.Getenv("DISCORD_API_VERSION")
	DISCORD_CDN_BASE_URL = os.Getenv("BASE_DISCORD_CDN_URL")
	DISCORD_APPLICATION_ID = os.Getenv("DISCORD_APPLICATION_ID")
	DISCORD_TESTING_GUILD_ID = os.Getenv("TESTING_GUILD_ID")
	DISCORD_BOT_TOKEN = os.Getenv("DISCORD_BOT_TOKEN")
	POSTGRESQL_CONNECTION_STRING = "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@" + os.Getenv("POSTGRES_HOST") + ":" + os.Getenv("POSTGRES_PORT") + "/" + os.Getenv("POSTGRES_DB")

	log.Info("Finished loading environment variables successfully")
}
