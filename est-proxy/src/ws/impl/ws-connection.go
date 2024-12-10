package impl

import (
	"github.com/gorilla/websocket"
	"sync"
)

type ConnectionImpl struct {
	conn      *websocket.Conn
	writeLock sync.Mutex
}

func NewConnectionImpl(conn *websocket.Conn) *ConnectionImpl {
	return &ConnectionImpl{
		conn:      conn,
		writeLock: sync.Mutex{},
	}
}

func (conn *ConnectionImpl) ReadMessage() ([]byte, error) {
	_, message, err := conn.conn.ReadMessage()
	return message, err
}

func (conn *ConnectionImpl) WriteMessage(message []byte) error {
	conn.writeLock.Lock()
	defer conn.writeLock.Unlock()

	return conn.conn.WriteMessage(websocket.TextMessage, message)
}

func (conn *ConnectionImpl) Close() error {
	return conn.conn.Close()
}
