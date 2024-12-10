package listener

import (
	"est-proxy/src/service"
	"est-proxy/src/service/impl"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type WsFigureListener struct {
	figureService service.WsFigureService
}

func NewWsFigureListener(figureService service.WsFigureService) *WsFigureListener {
	return &WsFigureListener{figureService: figureService}
}

func (l *WsFigureListener) Listen(ctx echo.Context) error {
	sessionUserId := impl.GetSessionUserId(ctx)
	if sessionUserId == nil {
		return ctx.String(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	boardId, err := uuid.Parse(ctx.QueryParam("boardId"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорретный запрос: %v", err))
	}

	statusError := l.figureService.Listen(ctx.Response().Writer, ctx.Request(), *sessionUserId, boardId)
	if statusError != nil {
		return statusError.Send(ctx)
	}

	return nil
}
