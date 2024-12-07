package listener

import (
	"est-proxy/src/ws"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
)

type MessageType int

const (
	ADD_FIGURE    MessageType = 0
	CHANGE_FIGURE MessageType = 1
	DELETE_FIGURE MessageType = 2
)

type WsFigureListener struct {
	chanel *ws.Channel
}

func NewWsFigureListener(chanel *ws.Channel) *WsFigureListener {
	return &WsFigureListener{chanel: chanel}
}

func (l *WsFigureListener) Listen(ctx echo.Context) error {
	// auth ?

	l.chanel.Listen(ctx.Response().Writer, ctx.Request(),
		func(boardId uuid.UUID, message []byte, conn *ws.Connection) {
			messageType, rawFigure := l.parseMessage(message)
			log.Printf("Received message: %v, figure: %s", messageType, rawFigure)

			// send to back?

			// send to rabbit?

			allConnections := l.getConnectionsForBoard(boardId)
			for _, c := range allConnections {
				if c == conn {
					continue
				}

				err := c.WriteMessage(message)
				if err != nil {
					log.Printf("Error writing message to connection: %v", err)
				}
			}
		})

	return nil
}

func (l *WsFigureListener) parseMessage(message []byte) (messageType MessageType, figure []byte) {
	messageType = MessageType(message[0])
	figure = message[1:]
	return
}

func (l *WsFigureListener) getConnectionsForBoard(boardId uuid.UUID) []*ws.Connection {
	return l.chanel.GetConnectionsForBoard(boardId)
}
