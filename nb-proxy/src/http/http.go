package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Listener struct {
}

func NewListener() *Listener {
	return &Listener{}
}

func (h *Listener) Serve() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/actuator", h.Actuator)

	e.Logger.Fatal(e.Start(":80"))
}

func (h *Listener) Actuator(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
