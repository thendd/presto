package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/coder/websocket"
)

type DiscordGatewayEventPayload struct {
	Opcode         string `json:"op"`
	Data           any    `json:"d"`
	SequenceNumber int    `json:"s"`
	Name           string `json:"t"`
}

func ReadIncomingMessages(connection *websocket.Conn, incomingMessages chan DiscordGatewayEventPayload) {
	for {
		_, rawMessage, err := connection.Read(context.Background())
		if err != nil {
			log.Println("Error while reading message: ", err)
		}

		var message DiscordGatewayEventPayload
		json.Unmarshal(rawMessage, &message)

		incomingMessages <- message
	}
}
