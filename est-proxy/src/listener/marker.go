package listener

import (
	"est-proxy/src/service"
	"est-proxy/src/service/impl"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type WsMarkerListener struct {
	markerService service.WsMarkerService
}

func NewWsMarkerListener(markerService service.WsMarkerService) *WsMarkerListener {
	return &WsMarkerListener{markerService: markerService}
}

func (l *WsMarkerListener) Listen(ctx echo.Context) error {
	sessionUserId := impl.GetSessionUserId(ctx)
	if sessionUserId == nil {
		return ctx.String(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	boardId, err := uuid.Parse(ctx.QueryParam("boardId"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорретный запрос: %v", err))
	}

	statusError := l.markerService.Listen(ctx.Response().Writer, ctx.Request(), *sessionUserId, boardId)
	if statusError != nil {
		return statusError.Send(ctx)
	}

	return nil
}
