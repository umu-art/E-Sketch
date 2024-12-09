package ws_channel

import (
	"est-proxy/src/ws/ws_channel/impl"
	"est-proxy/src/ws/ws_connection"
	"github.com/google/uuid"
	"net/http"
)

type Channel interface {
	Listen(responseWriter http.ResponseWriter, request *http.Request, onMessage impl.HandlerFunc)
	GetConnectionsForBoard(boardId uuid.UUID) []ws_connection.Connection
}

func NewChannel() Channel {
	return impl.NewChannelImpl()
}
