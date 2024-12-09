package impl

import (
	"est-proxy/src/ws/ws_connection"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
	"net/http"
)

type ChannelImpl struct {
	upgrader    *websocket.Upgrader
	connections *ConnectionsMap
}

func NewChannelImpl() *ChannelImpl {
	return &ChannelImpl{
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

type HandlerFunc func(boardId uuid.UUID, message []byte, conn ws_connection.Connection)

type AuthFunc func(boardId uuid.UUID, userId uuid.UUID) bool

func (channel *ChannelImpl) Listen(responseWriter http.ResponseWriter, request *http.Request, onMessage HandlerFunc) {
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

	boardId, err := uuid.Parse(request.URL.Query().Get("boardId"))
	if err != nil {
		log.Printf("Failed to parse boardId: %v", err)
		return
	}

	connection := ws_connection.NewConnection(ws)
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

func (channel *ChannelImpl) GetConnectionsForBoard(boardId uuid.UUID) []ws_connection.Connection {
	return channel.connections.GetConnections(boardId)
}
