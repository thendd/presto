package ws

import (
	"encoding/json"
	"io"
	"presto/internal/log"

	"github.com/gorilla/websocket"
)

type EventPayload struct {
	Opcode         string `json:"op"`
	Data           any    `json:"d"`
	SequenceNumber int    `json:"s"`
	Name           string `json:"t"`
}

func OnEventReceive(incomingEvents chan EventPayload, connection *websocket.Conn) {
	log.Info("Started listening to WSS' messages")
	for {
		_, reader, err := connection.NextReader()
		if err != nil {
			log.Error("Error while reading event: ", err)
		}

		rawEvent, err := io.ReadAll(reader)
		var event EventPayload
		json.Unmarshal(rawEvent, &event)

		incomingEvents <- event
	}
}
