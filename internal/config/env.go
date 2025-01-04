package config

import (
	"log"

	database_config "presto/internal/database/config"
	discord_config "presto/internal/discord/config"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() {
	log.Println("Started loading environment variables")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	discord_config.Load()
	database_config.Load()

	log.Println("Finished loading environment variables successfully")
}
