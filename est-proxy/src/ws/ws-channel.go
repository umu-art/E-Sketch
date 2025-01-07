package ws

import (
	"github.com/google/uuid"
	"net/http"
)

type Channel interface {
	Listen(responseWriter http.ResponseWriter, request *http.Request, onMessage HandlerFunc)
	GetConnectionsForBoard(boardId uuid.UUID) []Connection
}

type HandlerFunc func(message []byte, conn Connection)
