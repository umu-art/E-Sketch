package ws

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	conn      *websocket.Conn
	writeLock sync.Mutex
}

func NewConnection(conn *websocket.Conn) *Connection {
	return &Connection{
		conn:      conn,
		writeLock: sync.Mutex{},
	}
}

func (conn *Connection) ReadMessage() ([]byte, error) {
	_, message, err := conn.conn.ReadMessage()
	return message, err
}

func (conn *Connection) WriteMessage(message []byte) error {
	conn.writeLock.Lock()
	defer conn.writeLock.Unlock()

	return conn.conn.WriteMessage(websocket.TextMessage, message)
}

func (conn *Connection) Close() error {
	return conn.conn.Close()
}
