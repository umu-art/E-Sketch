package http

import (
	"est-proxy/src/listener"
	"est-proxy/src/service/impl"
	"est_proxy_go/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.elastic.co/apm/module/apmechov4/v2"
	"net/http"
	"strings"
)

type Listener struct {
	boardListener  handlers.BoardAPI
	userListener   handlers.UserAPI
	figureListener *listener.WsFigureListener
	markerListener *listener.WsMarkerListener
	gptListener    *listener.GptListener
}

func NewListener(
	boardListener handlers.BoardAPI,
	userListener handlers.UserAPI,
	figureListener *listener.WsFigureListener,
	markerListener *listener.WsMarkerListener,
	gptListener *listener.GptListener) *Listener {

	return &Listener{
		boardListener:  boardListener,
		userListener:   userListener,
		figureListener: figureListener,
		markerListener: markerListener,
		gptListener:    gptListener,
	}
}

func (h *Listener) Serve() {

	e := echo.New()

	e.Use(middleware.Recover())

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "/actuator")
		},
	}))

	e.Use(apmechov4.Middleware(
		apmechov4.WithRequestIgnorer(
			func(req *http.Request) bool {
				return strings.Contains(req.URL.Path, "/actuator")
			},
		),
	))

	e.Use(impl.SessionMiddleware)

	e.GET("/actuator", h.Actuator)

	e.Any("/proxy/figure/ws", h.figureListener.Listen)
	e.Any("/proxy/marker/ws", h.markerListener.Listen)

	handlers.RouteBoardAPI(e, h.boardListener)
	handlers.RouteUserAPI(e, h.userListener)
	handlers.RouteGptAPI(e, h.gptListener)

	e.Logger.Fatal(e.Start(":8080"))
}

func (h *Listener) Actuator(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
