package main

import (
	"os"
	"os/signal"
	"presto/internal/bot"
	"presto/internal/bot/commands"
	"presto/internal/config"
	"presto/internal/database"
)

func main() {
	config.LoadEnvironmentVariables()
	database.Connect()

	var localApplicationCommands = []bot.ApplicationCommandWithHandler{
		commands.Ping,
		commands.WarnUserCommand,
		commands.WarnSlashCommand,
		commands.WarnMessageCommand,
		commands.GuildSettings,
		commands.Clear,
	}

	session := bot.NewSession(localApplicationCommands)
	err := session.Open()
	if err != nil {
		panic(err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch
}
