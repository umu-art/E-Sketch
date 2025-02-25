package errors

import (
	"errors"
	"github.com/labstack/echo/v4"
)

type StatusError struct {
	HttpStatus int
	message    string
}

func NewStatusError(httpStatus int, message string) *StatusError {
	return &StatusError{HttpStatus: httpStatus, message: message}
}

func (e *StatusError) Send(ctx echo.Context) error {
	return ctx.String(e.HttpStatus, e.message)
}

func (e *StatusError) GetMessage() string {
	return e.message
}

var ErrRabbitChannelClosed = errors.New("channel is closed")
