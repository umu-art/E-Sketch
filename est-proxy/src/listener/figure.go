package listener

import (
	"est-proxy/src/service"
	"github.com/labstack/echo/v4"
	"log"
)

type WsFigureListener struct {
	figureService service.WsFigureService
}

func NewWsFigureListener(figureService service.WsFigureService) *WsFigureListener {
	return &WsFigureListener{figureService: figureService}
}

func (l *WsFigureListener) Listen(ctx echo.Context) error {
	err := l.figureService.Listen(ctx.Response().Writer, ctx.Request())
	if err != nil {
		log.Printf("error listening on figure service: %v", err)
	}
	return nil
}
