package service

import (
	"context"
	"est-proxy/src/errors"
	proxymodels "est_proxy_go/models"
	"github.com/google/uuid"
)

type BoardService interface {
	GetByUuid(ctx context.Context, userId *uuid.UUID, boardId *uuid.UUID) (*proxymodels.BoardDto, *errors.StatusError)
	List(ctx context.Context, userId *uuid.UUID) (*proxymodels.BoardListDto, *errors.StatusError)
	Create(ctx context.Context, userId *uuid.UUID, createRequest *proxymodels.CreateRequest) (*proxymodels.BoardDto, *errors.StatusError)
	Update(ctx context.Context, userId *uuid.UUID, boardId *uuid.UUID, updateRequest *proxymodels.CreateRequest) (*proxymodels.BoardDto, *errors.StatusError)
	DeleteBoard(ctx context.Context, userId *uuid.UUID, boardId *uuid.UUID) *errors.StatusError
	Share(ctx context.Context, userId *uuid.UUID, boardId *uuid.UUID, shareBoardDto *proxymodels.ShareBoardDto) *errors.StatusError
	ChangeAccess(ctx context.Context, userId *uuid.UUID, boardId *uuid.UUID, shareBoardDto *proxymodels.ShareBoardDto) *errors.StatusError
	Unshare(ctx context.Context, userId *uuid.UUID, boardId *uuid.UUID, unshareRequest *proxymodels.UnshareRequest) *errors.StatusError
}
