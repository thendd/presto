package config

import (
	"log"

	"presto/internal/discord/config"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() {
	log.Println("Started loading environment variables")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	config.LoadDiscordConfig()

	log.Println("Finished loading environment variables successfully")
}
