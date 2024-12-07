package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
	"net/http"
)

type Channel struct {
	upgrader    *websocket.Upgrader
	connections *ConnectionsMap
}

func NewChannel() *Channel {
	return &Channel{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // TODO: Check origin
			},
		},
		connections: NewConnectionsMap(),
	}
}

type HandlerFunc func(boardId uuid.UUID, message []byte, conn *Connection)

type AuthFunc func(boardId uuid.UUID, userId uuid.UUID) bool

func (channel *Channel) Listen(responseWriter http.ResponseWriter, request *http.Request, onMessage HandlerFunc) {
	ws, err := channel.upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
	}

	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Printf("Failed to close connection: %v", err)
		}
	}(ws)

	boardId := uuid.UUID{} // Pass in headers?

	connection := NewConnection(ws)
	channel.connections.Save(boardId, connection)
	defer channel.connections.Remove(boardId, connection)

	for {
		message, err := connection.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		onMessage(boardId, message, connection)
	}
}

func (channel *Channel) GetConnectionsForBoard(boardId uuid.UUID) []*Connection {
	return channel.connections.GetConnections(boardId)
}
