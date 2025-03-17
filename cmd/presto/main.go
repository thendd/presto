package main

import (
	"flag"
	"os"
	"os/signal"
	"presto/internal/bot"
	"presto/internal/bot/commands"
	"presto/internal/config"
	"presto/internal/database"
	"presto/internal/log"
)

func init() {
	debugFlag := flag.Uint("debug_flags", 0, "Set the debug flags")
	flag.Parse()

	if *debugFlag == 0 {
		panic("debug_flags not set")
	}

	if err := log.New(os.Stdout, log.DebugLevel(*debugFlag)); err != nil {
		panic(err)
	}
}

func main() {
	config.LoadEnvironmentVariables()
	database.Connect()

	var localApplicationCommands = []bot.ApplicationCommandWithHandler{
		commands.Ping,
		commands.WarnUserCommand,
		commands.WarnSlashCommand,
		commands.WarnMessageCommand,
		commands.GuildSettings,
	}

	session := bot.NewSession(localApplicationCommands)
	err := session.Open()
	if err != nil {
		log.Fatal(err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch
}
