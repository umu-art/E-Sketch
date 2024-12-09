package impl

import (
	"context"
	"est-proxy/src/ws/ws_channel"
	"est-proxy/src/ws/ws_connection"
	estbackapi "est_back_go"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type MessageType int

const (
	ADD_FIGURE    MessageType = 0
	CHANGE_FIGURE MessageType = 1
	DELETE_FIGURE MessageType = 2
)

type WsFigureServiceImpl struct {
	channel   ws_channel.Channel
	figureApi *estbackapi.FigureAPIService
}

func NewWsServiceListenerImpl(channel ws_channel.Channel, figureApi *estbackapi.FigureAPIService) *WsFigureServiceImpl {
	return &WsFigureServiceImpl{channel: channel, figureApi: figureApi}
}

func (l *WsFigureServiceImpl) Listen(writer http.ResponseWriter, request *http.Request) error {
	l.channel.Listen(writer, request,
		func(boardId uuid.UUID, message []byte, conn ws_connection.Connection) {
			messageType, rawFigure := l.parseMessage(message)
			log.Printf("Received message: %v, figure: %s", messageType, rawFigure)

			switch messageType {
			case ADD_FIGURE:
				figureId, _ := l.addFigure(request.Context(), &boardId)
				if figureId.String() == "" {
					log.Printf("Figure Id is empty")
				}
				//TODO send it to front
			case CHANGE_FIGURE:
				figureId, err := uuid.Parse(request.URL.Query().Get("boardId"))
				if err != nil {
					log.Printf("Failed to parse figureId: %v", err)
					return
				}
				err = l.deleteFigure(request.Context(), &figureId)
				if err != nil {
					log.Printf("Failed to delete figure: %v", err)
				}

			case DELETE_FIGURE:
				figureId, err := l.getFigureId(rawFigure)
				if err != nil {
					log.Printf("Failed to parse figureId: %v", err)
					return
				}
				err = l.deleteFigure(request.Context(), figureId)
			}

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

func (l *WsFigureServiceImpl) parseMessage(message []byte) (messageType MessageType, figure []byte) {
	messageType = MessageType(message[0])
	figure = message[1:]
	return
}

func (l *WsFigureServiceImpl) getConnectionsForBoard(boardId uuid.UUID) []ws_connection.Connection {
	return l.channel.GetConnectionsForBoard(boardId)
}

func (l *WsFigureServiceImpl) addFigure(ctx context.Context, boardId *uuid.UUID) (*uuid.UUID, error) {
	figureIdDto, _, err := l.figureApi.CreateFigure(ctx, boardId.String()).Execute()
	if err != nil {
		log.Printf("Error creating figure: %v", err)
		return nil, err
	}

	figureId, err := uuid.Parse(figureIdDto.GetId())
	if err != nil {
		log.Printf("Error parsing figure id: %v", err)
	}

	return &figureId, nil
}

func (l *WsFigureServiceImpl) deleteFigure(ctx context.Context, figureId *uuid.UUID) error {
	_, err := l.figureApi.DeleteFigure(ctx, figureId.String()).Execute()
	if err != nil {
		log.Printf("Error creating figure: %v", err)
		return err
	}

	return nil
}

func (l *WsFigureServiceImpl) changeFigure(ctx context.Context, figureId *uuid.UUID, rawFigure []byte) error {
	figureDto := estbackapi.FigureDto{} //TODO parse rawFigure
	_, err := l.figureApi.UpdateFigure(ctx, figureId.String()).FigureDto(figureDto).Execute()
	if err != nil {
		log.Printf("Error creating figure: %v", err)
		return err
	}

	return nil
}

func (l *WsFigureServiceImpl) getFigureId(rawFigure []byte) (*uuid.UUID, error) {
	idStr := string(rawFigure[1:37])

	parsedUUID, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse UUID: " + err.Error())
	}

	return &parsedUUID, nil
}
