package http

import (
	"est-proxy/src/listener"
	"est-proxy/src/service/impl"
	"est_proxy_go/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.elastic.co/apm/module/apmechov4/v2"
	"net/http"
)

type Listener struct {
	boardListener  handlers.BoardAPI
	userListener   handlers.UserAPI
	figureListener *listener.WsFigureListener
}

func NewListener(
	boardListener handlers.BoardAPI,
	userListener handlers.UserAPI,
	figureListener *listener.WsFigureListener) *Listener {

	return &Listener{
		boardListener:  boardListener,
		userListener:   userListener,
		figureListener: figureListener,
	}
}

func (h *Listener) Serve() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(apmechov4.Middleware())
	e.Use(impl.SessionMiddleware)

	e.GET("/actuator", h.Actuator)

	e.Any("/proxy/ws", h.figureListener.Listen)

	handlers.RouteBoardAPI(e, h.boardListener)
	handlers.RouteUserAPI(e, h.userListener)

	e.Logger.Fatal(e.Start(":8080"))
}

func (h *Listener) Actuator(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
