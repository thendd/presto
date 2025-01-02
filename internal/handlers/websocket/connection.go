package ws

import (
	"context"
	"log"

	"github.com/coder/websocket"
)

var Connection *websocket.Conn

func ConnnectToWebsocket(URL string) *websocket.Conn {
	ws, _, err := websocket.Dial(context.Background(), URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	Connection = ws

	return Connection
}
