package main

import (
	"os"
	"os/signal"
	"presto/internal/config"
	"presto/internal/database"
	ws "presto/internal/discord/websocket"
)

func main() {
	config.LoadEnvironmentVariables()
	database.Connect()

	session := ws.NewSession()
	err := session.Open()
	if err != nil {
		panic(err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch
}
