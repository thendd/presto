package ws

import (
	"context"
	"fmt"
	"log"

	"presto/configs"

	"github.com/coder/websocket"
)

func Identify(connection *websocket.Conn) {
	// TODO: currently using "32767" to activate all intents. Figure out a better approach later.
	toSend := fmt.Sprintf("{\"op\":2,\"d\":{\"token\":\"%s\",\"properties\":{\"os\":null,\"browser\":null,\"device\":null},\"intents\":32767}}", configs.DISCORD_BOT_TOKEN)

	err := connection.Write(context.Background(), websocket.MessageText, []byte(toSend))
	if err != nil {
		log.Fatal(err)
	}
}
