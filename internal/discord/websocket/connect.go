package ws

import (
	"net/http"
	"presto/internal/log"

	"github.com/gorilla/websocket"
)

var Connection *websocket.Conn

func Connnect(URL string) *websocket.Conn {
	log.Info("Started connection to Discord's WSS")

	dialer := websocket.Dialer{}
	ws, _, err := dialer.Dial(URL, http.Header{})
	if err != nil {
		log.Fatal(err.Error())
	}

	Connection = ws

	log.Info("Performed websocket handshake with Discord's WSS")

	return ws
}
