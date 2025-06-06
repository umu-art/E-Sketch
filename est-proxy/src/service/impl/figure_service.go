package impl

import (
	"context"
	"encoding/json"
	"errors"
	esterrors "est-proxy/src/errors"
	"est-proxy/src/repository"
	"est-proxy/src/utils"
	"est-proxy/src/ws"
	estbackapi "est_back_go"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type MessageType int

const (
	AddFigure    MessageType = 0
	ChangeFigure MessageType = 1
	DeleteFigure MessageType = 2
	UpdateFigure MessageType = 3
	GetFigures   MessageType = 4
)

type WsFigureServiceImpl struct {
	channel      ws.Channel
	figureApi    *estbackapi.FigureAPIService
	boardApi     *estbackapi.BoardAPIService
	topic        repository.RabbitTopic
	figureBuffer *utils.FigureBuffer
}

func NewWsFigureServiceImpl(
	channel ws.Channel,
	figureApi *estbackapi.FigureAPIService,
	boardApi *estbackapi.BoardAPIService,
	topic repository.RabbitTopic,
	figureBuffer *utils.FigureBuffer,
) *WsFigureServiceImpl {

	service := &WsFigureServiceImpl{
		channel:      channel,
		figureApi:    figureApi,
		boardApi:     boardApi,
		topic:        topic,
		figureBuffer: figureBuffer,
	}

	err := topic.Subscribe(func(message []byte) {
		service.handleChangedFigure(message)
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	return service
}

func (l *WsFigureServiceImpl) Listen(writer http.ResponseWriter, request *http.Request, userId uuid.UUID, boardId uuid.UUID) *esterrors.StatusError {
	if !l.checkAvailability(userId, boardId) {
		return esterrors.NewStatusError(http.StatusForbidden, "Недостаточно прав")
	}

	l.channel.Listen(writer, request,
		func(message []byte, conn ws.Connection) {
			messageType, rawFigure := l.parseMessage(message)

			switch messageType {
			case AddFigure:
				err := l.addFigure(conn, &boardId)
				if err != nil {
					log.Printf("Failed to add figure: %v", err)
					return
				}
				break
			case ChangeFigure:
				figureId, err := l.getFigureId(rawFigure)
				if err != nil {
					log.Printf("Failed to parse figureId: %v", err)
					return
				}
				err = l.changeFigure(boardId, figureId, rawFigure)
				if err != nil {
					log.Printf("Failed to change figure: %v", err)
					return
				}
				break
			case DeleteFigure:
				figureId, err := l.getFigureId(rawFigure)
				if err != nil {
					log.Printf("Failed to parse figureId: %v", err)
					return
				}
				err = l.deleteFigure(boardId, figureId)
				if err != nil {
					log.Printf("Failed to delete figure: %v", err)
					return
				}
				break
			case UpdateFigure:
				figureId, err := l.getFigureId(rawFigure)
				if err != nil {
					log.Printf("Failed to parse figureId: %v", err)
					return
				}
				err = l.updateFigure(boardId, figureId, rawFigure)
				if err != nil {
					log.Printf("Failed to update figure: %v", err)
					return
				}
				break
			case GetFigures:
				err := l.getFigures(conn, &boardId)
				if err != nil {
					log.Printf("Failed to get board figures: %v", err)
					return
				}
				break
			}
		})

	return nil
}

func (l *WsFigureServiceImpl) parseMessage(message []byte) (messageType MessageType, figure []byte) {
	messageType = MessageType(message[0])
	figure = message[1:]
	return
}

func (l *WsFigureServiceImpl) getConnectionsForBoard(boardId uuid.UUID) []ws.Connection {
	return l.channel.GetConnectionsForBoard(boardId)
}

func (l *WsFigureServiceImpl) addFigure(conn ws.Connection, boardId *uuid.UUID) error {
	figureIdDto, _, err := l.figureApi.CreateFigure(context.Background(), boardId.String()).Execute()
	if err != nil {
		return fmt.Errorf("error back creating figure: %v", err)
	}

	figureId, err := uuid.Parse(figureIdDto.GetId())
	if err != nil {
		return fmt.Errorf("error parsing figureId: %v", err)
	}

	err = conn.WriteMessage([]byte(figureId.String()))
	if err != nil {
		return fmt.Errorf("error writing message to connection: %v", err)
	}

	return nil
}

func (l *WsFigureServiceImpl) deleteFigure(boardId uuid.UUID, figureId uuid.UUID) error {
	l.notifyChangedFigure(boardId, []byte("-"+figureId.String()))

	if l.figureBuffer.Remove(figureId.String()) {
		return nil
	}

	_, err := l.figureApi.DeleteFigure(context.Background(), figureId.String()).Execute()
	if err != nil {
		return fmt.Errorf("error back deleting figure: %v", err)
	}

	return nil
}

func (l *WsFigureServiceImpl) changeFigure(boardId uuid.UUID, figureId uuid.UUID, rawFigure []byte) error {
	l.notifyChangedFigure(boardId, rawFigure)

	figureDto := estbackapi.FigureDto{
		Data: string(rawFigure),
	}

	_, err := l.figureApi.UpdateFigure(context.Background(), figureId.String()).FigureDto(figureDto).Execute()
	if err != nil {
		return fmt.Errorf("error back updating figure: %v", err)
	}

	return nil
}

func (l *WsFigureServiceImpl) updateFigure(boardId uuid.UUID, figureId uuid.UUID, rawFigure []byte) error {
	l.notifyChangedFigure(boardId, []byte("+"+string(rawFigure)))

	l.figureBuffer.Add(figureId.String(), rawFigure[37:])

	return nil
}

func (l *WsFigureServiceImpl) getFigures(conn ws.Connection, boardId *uuid.UUID) error {
	figureListDto, _, err := l.figureApi.ListByBoardId(context.Background(), boardId.String()).Execute()
	if err != nil {
		return fmt.Errorf("error back getting figures: %v", err)
	}
	for _, figureDto := range figureListDto.Figures {
		err = conn.WriteMessage([]byte(figureDto.Data))
		if err != nil {
			log.Printf("error writing message to connection: %v", err)
		}
	}
	return nil
}

func (l *WsFigureServiceImpl) UpdateFigure(figureId string, figureData []byte) {
	figureDto, _, err := l.figureApi.GetFigure(context.Background(), figureId).Execute()
	if err != nil {
		log.Printf("Error getting figure: %v", err)
	}

	figureDto.Data += string(figureData)

	_, err = l.figureApi.UpdateFigure(context.Background(), figureId).FigureDto(*figureDto).Execute()
	if err != nil {
		log.Printf("Error updating figure: %v", err)
	}
}

type FigureChangeMessage struct {
	BoardId    uuid.UUID `json:"board-id"`
	FigureData []byte    `json:"figure-data"`
}

func (l *WsFigureServiceImpl) notifyChangedFigure(boardId uuid.UUID, figureData []byte) {
	message := FigureChangeMessage{
		BoardId:    boardId,
		FigureData: figureData,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal json ws message: %v", err)
		return
	}

	err = l.topic.Publish(jsonMessage)
	if err != nil && !errors.Is(err, esterrors.ErrRabbitChannelClosed) {
		log.Printf("Failed to publish json ws message: %v", err)
		return
	}
}

func (l *WsFigureServiceImpl) handleChangedFigure(message []byte) {
	var figureChangeMessage FigureChangeMessage
	err := json.Unmarshal(message, &figureChangeMessage)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}

	connections := l.getConnectionsForBoard(figureChangeMessage.BoardId)
	for _, connection := range connections {
		err := connection.WriteMessage(figureChangeMessage.FigureData)
		if err != nil {
			log.Printf("Error writing message: %v", err)
			continue
		}
	}
}

func (l *WsFigureServiceImpl) getFigureId(rawFigure []byte) (uuid.UUID, error) {
	idStr := string(rawFigure[1:37])

	parsedUUID, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to parse UUID: " + err.Error())
	}

	return parsedUUID, nil
}

func (l *WsFigureServiceImpl) checkAvailability(userId uuid.UUID, boardId uuid.UUID) bool {
	boardDto, _, err := l.boardApi.GetByUuid(context.Background(), boardId.String()).Execute()
	if err != nil {
		log.Printf("Error getting board: %v", err)
		return false
	}

	return utils.GetAccessLevel(&userId, boardDto) != utils.NONE
}
