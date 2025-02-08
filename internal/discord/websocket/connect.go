package ws

import (
	"context"
	"presto/internal/log"

	"github.com/coder/websocket"
)

var Connection *websocket.Conn

func Connnect(URL string) *websocket.Conn {
	log.Info("Started connection to Discord's WSS")

	ws, _, err := websocket.Dial(context.Background(), URL, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	Connection = ws

	log.Info("Performed websocket handshake with Discord's WSS")

	return Connection
}
