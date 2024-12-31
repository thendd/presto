package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/coder/websocket"
)

func ReadIncomingMessages(connection *websocket.Conn, incomingMessages chan map[string]any) {
	for {
		_, rawMessage, err := connection.Read(context.Background())
		if err != nil {
			log.Println("Error while reading message: ", err)
		}

		var message map[string]any
		json.Unmarshal(rawMessage, &message)

		incomingMessages <- message
	}
}
