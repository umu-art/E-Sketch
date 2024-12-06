package service

import (
	"context"
	"est-proxy/src/errors"
	"est-proxy/src/repository/user_repository"
	"est-proxy/src/service/impl"
	proxymodels "est_proxy_go/models"
	"github.com/google/uuid"
)

type UserService interface {
	GetUserById(ctx context.Context, userId *uuid.UUID) (*proxymodels.UserDto, *errors.StatusError)
	Login(ctx context.Context, authDto *proxymodels.AuthDto) (*string, *errors.StatusError)
	Register(ctx context.Context, registerDto *proxymodels.RegisterDto) (*string, *errors.StatusError)
	Search(ctx context.Context, query string) (*[]proxymodels.UserDto, *errors.StatusError)
}

func NewUserService(userRepository user_repository.UserRepository) UserService {
	return impl.NewUserServiceImpl(userRepository)
}
