package ws

import (
	"context"
	"log"

	"github.com/coder/websocket"
)

var Connection *websocket.Conn

func Connnect(URL string) *websocket.Conn {
	log.Println("Started connection to Discord's WSS")

	ws, _, err := websocket.Dial(context.Background(), URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	Connection = ws

	log.Println("Performed websocket handshake with Discord's WSS")

	return Connection
}
