package ws

import (
	"context"
	"encoding/json"
	"presto/internal/log"
)

type EventPayload struct {
	Opcode         string `json:"op"`
	Data           any    `json:"d"`
	SequenceNumber int    `json:"s"`
	Name           string `json:"t"`
}

func OnEventReceive(incomingEvents chan EventPayload) {
	log.Info("Started listening to WSS' messages")
	for {
		_, rawEvent, err := Connection.Read(context.Background())
		if err != nil {
			log.Error("Error while reading event: ", err)
		}

		var event EventPayload
		json.Unmarshal(rawEvent, &event)

		incomingEvents <- event
	}
}
