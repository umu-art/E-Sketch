package api

import (
	"context"
	"est-proxy/src/errors"
)

type PreviewApi interface {
	GetTokens(boardIds []string, ctx context.Context) (map[string]string, *errors.StatusError)
	GetPreview(
		boardId string,
		width float32, height float32,
		xLeft float32, yTop float32,
		xRight float32, yBottom float32,
		ctx context.Context,
	) ([]byte, *errors.StatusError)
}
