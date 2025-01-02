package ws

import (
	"context"
	"encoding/json"
	"log"
)

type DiscordGatewayEventPayload struct {
	Opcode         string `json:"op"`
	Data           any    `json:"d"`
	SequenceNumber int    `json:"s"`
	Name           string `json:"t"`
}

func ReadIncomingMessages(incomingMessages chan DiscordGatewayEventPayload) {
	for {
		_, rawMessage, err := Connection.Read(context.Background())
		if err != nil {
			log.Println("Error while reading message: ", err)
		}

		var message DiscordGatewayEventPayload
		json.Unmarshal(rawMessage, &message)

		incomingMessages <- message
	}
}
