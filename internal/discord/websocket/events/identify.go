package events

import (
	"fmt"
	"presto/internal/log"

	"presto/internal/discord/config"

	"github.com/gorilla/websocket"
)

func SendIdentify(connection *websocket.Conn) {
	log.Info("Triggered initial handshake")

	// Currently using only `GUILDS` intents
	toSend := fmt.Sprintf("{\"op\":2,\"d\":{\"token\":\"%s\",\"properties\":{\"os\":null,\"browser\":null,\"device\":null},\"intents\":1}}", config.DISCORD_BOT_TOKEN)

	err := connection.WriteMessage(websocket.TextMessage, []byte(toSend))
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Handshake ocurred successfully")
}
