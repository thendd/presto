package config

import (
	"presto/internal/log"

	database_config "presto/internal/database/config"
	discord_config "presto/internal/discord/config"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() {
	log.Info("Started loading environment variables")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
	discord_config.Load()
	database_config.Load()

	log.Info("Finished loading environment variables successfully")
}
