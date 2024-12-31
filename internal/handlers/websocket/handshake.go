package ws

import (
	"context"
	"fmt"
	"log"

	"presto/configs"

	"github.com/coder/websocket"
)

func SendIdentify(connection *websocket.Conn) {
	// Currently using only `GUILDS` intents
	toSend := fmt.Sprintf("{\"op\":2,\"d\":{\"token\":\"%s\",\"properties\":{\"os\":null,\"browser\":null,\"device\":null},\"intents\":1}}", configs.DISCORD_BOT_TOKEN)

	err := connection.Write(context.Background(), websocket.MessageText, []byte(toSend))
	if err != nil {
		log.Fatal(err)
	}
}
