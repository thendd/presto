package events

import (
	"fmt"
	"presto/internal/log"

	"github.com/gorilla/websocket"
)

// Sends a heartbeat to Discord's WSS
func SendHeartbeat(lastSequenceNumber any, connection *websocket.Conn) {
	log.Info("Started sending heartbeat")
	var heartbeat string

	if lastSequenceNumber != nil {
		heartbeat = fmt.Sprintf("{\"op\": 1, \"d\": %v}", lastSequenceNumber)
	} else {
		heartbeat = "{\"op\": 1, \"d\": null}"
	}

	err := connection.WriteMessage(websocket.TextMessage, []byte(heartbeat))
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Sent heartbeat")
}
