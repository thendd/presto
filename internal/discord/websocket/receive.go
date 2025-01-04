package ws

import (
	"context"
	"encoding/json"
	"log"
)

type EventPayload struct {
	Opcode         string `json:"op"`
	Data           any    `json:"d"`
	SequenceNumber int    `json:"s"`
	Name           string `json:"t"`
}

func OnEventReceive(incomingEvents chan EventPayload) {
	log.Println("Started listening to WSS' messages")
	for {
		_, rawEvent, err := Connection.Read(context.Background())
		if err != nil {
			log.Println("Error while reading event: ", err)
		}

		log.Println(string(rawEvent))

		var event EventPayload
		json.Unmarshal(rawEvent, &event)

		incomingEvents <- event
	}
}
