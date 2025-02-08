package events

import (
	"context"
	"fmt"
	"presto/internal/log"

	ws "presto/internal/discord/websocket"

	"github.com/coder/websocket"
)

// Sends a heartbeat to Discord's WSS
func SendHeartbeat(lastSequenceNumber any) {
	log.Info("Started sending heartbeat")
	var heartbeat string

	if lastSequenceNumber != nil {
		heartbeat = fmt.Sprintf("{\"op\": 1, \"d\": %v}", lastSequenceNumber)
	} else {
		heartbeat = "{\"op\": 1, \"d\": null}"
	}

	err := ws.Connection.Write(context.Background(), websocket.MessageText, []byte(heartbeat))
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Sent heartbeat")
}
