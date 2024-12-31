package ws

import (
	"context"
	"log"

	"github.com/coder/websocket"
)

func ConnnectToWebsocket(URL string) *websocket.Conn {
	ws, _, err := websocket.Dial(context.Background(), URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	return ws
}
