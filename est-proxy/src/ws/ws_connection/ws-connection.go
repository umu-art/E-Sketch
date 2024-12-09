package ws_connection

import (
	"est-proxy/src/ws/ws_connection/impl"
	"github.com/gorilla/websocket"
)

type Connection interface {
	ReadMessage() ([]byte, error)
	WriteMessage(message []byte) error
	Close() error
}

func NewConnection(conn *websocket.Conn) Connection {
	return impl.NewConnectionImpl(conn)
}
