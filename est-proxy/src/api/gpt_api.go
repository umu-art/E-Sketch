package api

import (
	"context"
	"est-proxy/src/errors"
)

type GptApi interface {
	Request(prompt string, image []byte, ctx context.Context) (string, *errors.StatusError)
}
