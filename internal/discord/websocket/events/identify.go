package events

import (
	"context"
	"fmt"
	"log"

	"presto/internal/discord/config"
	ws "presto/internal/discord/websocket"

	"github.com/coder/websocket"
)

func SendIdentify() {
	log.Println("Triggered initial handshake")

	// Currently using only `GUILDS` intents
	toSend := fmt.Sprintf("{\"op\":2,\"d\":{\"token\":\"%s\",\"properties\":{\"os\":null,\"browser\":null,\"device\":null},\"intents\":1}}", config.DISCORD_BOT_TOKEN)

	err := ws.Connection.Write(context.Background(), websocket.MessageText, []byte(toSend))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Handshake ocurred successfully")
}
