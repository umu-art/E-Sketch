package impl

import (
	"context"
	"encoding/json"
	"est-proxy/src/errors"
	"est-proxy/src/repository"
	"est-proxy/src/service"
	"est-proxy/src/utils"
	"est-proxy/src/ws"
	estbackapi "est_back_go"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type WsMarkerServiceImpl struct {
	channel     ws.Channel
	boardApi    *estbackapi.BoardAPIService
	userService service.UserService
	topic       repository.Topic
}

func NewWsMarkerServiceImpl(
	channel ws.Channel,
	boardApi *estbackapi.BoardAPIService,
	userService service.UserService,
	topic repository.Topic,
) *WsMarkerServiceImpl {

	markerService := &WsMarkerServiceImpl{
		channel:     channel,
		boardApi:    boardApi,
		userService: userService,
		topic:       topic,
	}

	err := topic.Subscribe(func(message []byte) {
		markerService.handleChangedMarker(message)
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	return markerService
}

func (l *WsMarkerServiceImpl) Listen(writer http.ResponseWriter, request *http.Request, userId uuid.UUID, boardId uuid.UUID) *errors.StatusError {
	//if !l.checkAvailability(userId, boardId) {
	//	return errors.NewStatusError(http.StatusForbidden, "Недостаточно прав")
	//}

	userDto, err := l.userService.GetUserById(request.Context(), &userId)
	if err != nil {
		return err
	}

	l.channel.Listen(writer, request,
		func(message []byte, conn ws.Connection) {
			l.updateMarker(boardId, message, []byte(userDto.Username))
		})

	return nil
}

func (l *WsMarkerServiceImpl) getConnectionsForBoard(boardId uuid.UUID) []ws.Connection {
	return l.channel.GetConnectionsForBoard(boardId)
}

func (l *WsMarkerServiceImpl) updateMarker(boardId uuid.UUID, message []byte, username []byte) {
	r := append(message, username...)
	l.notifyChangedMarker(boardId, r)
}

type MarkerChangeMessage struct {
	BoardId    uuid.UUID `json:"board-id"`
	MarkerData []byte    `json:"marker-data"`
}

func (l *WsMarkerServiceImpl) notifyChangedMarker(boardId uuid.UUID, markerData []byte) {
	message := MarkerChangeMessage{
		BoardId:    boardId,
		MarkerData: markerData,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	err = l.topic.Publish(jsonMessage)
	if err != nil {
		log.Printf("Error publishing message: %v", err)
		return
	}
}

func (l *WsMarkerServiceImpl) handleChangedMarker(message []byte) {
	var markerChangeMessage MarkerChangeMessage
	err := json.Unmarshal(message, &markerChangeMessage)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}

	connections := l.getConnectionsForBoard(markerChangeMessage.BoardId)
	for _, connection := range connections {
		err := connection.WriteMessage(markerChangeMessage.MarkerData)
		if err != nil {
			log.Printf("Error writing message: %v", err)
			continue
		}
	}
}

func (l *WsMarkerServiceImpl) checkAvailability(userId uuid.UUID, boardId uuid.UUID) bool {
	boardDto, _, err := l.boardApi.GetByUuid(context.Background(), boardId.String()).Execute()
	if err != nil {
		log.Printf("Error getting board: %v", err)
		return false
	}

	return utils.GetAccessLevel(&userId, boardDto) != utils.NONE
}
