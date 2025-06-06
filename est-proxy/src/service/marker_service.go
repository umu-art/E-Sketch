package service

import (
	"est-proxy/src/errors"
	"github.com/google/uuid"
	"net/http"
)

type WsMarkerService interface {
	Listen(writer http.ResponseWriter, request *http.Request, userId uuid.UUID, boardId uuid.UUID) *errors.StatusError
}
