package errors

import "github.com/labstack/echo/v4"

type StatusError struct {
	httpStatus int
	message    string
}

func NewStatusError(httpStatus int, message string) *StatusError {
	return &StatusError{httpStatus: httpStatus, message: message}
}

func (e StatusError) Send(ctx echo.Context) error {
	return ctx.String(e.httpStatus, e.message)
}
