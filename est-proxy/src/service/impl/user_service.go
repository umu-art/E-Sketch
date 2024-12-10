package impl

import (
	"context"
	"est-proxy/src/errors"
	"est-proxy/src/mapper"
	"est-proxy/src/repository"
	"est-proxy/src/utils"
	proxymodels "est_proxy_go/models"
	"github.com/google/uuid"
	"net/http"
)

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserServiceImpl(userRepository repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{userRepository: userRepository}
}

func (u UserServiceImpl) GetUserById(ctx context.Context, userId *uuid.UUID) (*proxymodels.UserDto, *errors.StatusError) {
	user := u.userRepository.GetUserByID(ctx, userId)
	if user == nil {
		return nil, errors.NewStatusError(http.StatusBadRequest, "Пользователь не найден")
	}

	return mapper.MapUserToProxy(*user.Public()), nil
}

func (u UserServiceImpl) Login(ctx context.Context, authDto *proxymodels.AuthDto) (*string, *errors.StatusError) {
	user := u.userRepository.GetUserByEmail(ctx, authDto.Email)
	if user == nil {
		return nil, errors.NewStatusError(http.StatusBadRequest, "Отсутвует или некорректный адрес почты или пароль")
	}

	if authDto.PasswordHash != user.PasswordHash {
		return nil, errors.NewStatusError(http.StatusBadRequest, "Отсутвует или некорректный адрес почты или пароль")
	}

	token := utils.GenerateUserJWTString(&user.ID)
	if token == nil {
		return nil, errors.NewStatusError(http.StatusUnauthorized, "Не получилось начать сессию")
	}

	return token, nil
}

func (u UserServiceImpl) Register(ctx context.Context, registerDto *proxymodels.RegisterDto) (*string, *errors.StatusError) {
	exists := u.userRepository.UserExistsByUsernameOrEmail(ctx, registerDto.Username, registerDto.Email)
	if exists {
		return nil, errors.NewStatusError(http.StatusConflict, "Занято имя пользователя или адрес электронной почты")
	}

	userId := u.userRepository.Create(ctx, registerDto.Username, registerDto.Email, registerDto.PasswordHash)
	if userId == nil {
		return nil, errors.NewStatusError(http.StatusInternalServerError, "Не получилось создать аккаунт")
	}

	token := utils.GenerateUserJWTString(userId)
	if token == nil {
		return nil, errors.NewStatusError(http.StatusUnauthorized, "Не получилось начать сессию")
	}

	return token, nil
}

func (u UserServiceImpl) Search(ctx context.Context, query string) (*[]proxymodels.UserDto, *errors.StatusError) {
	users := u.userRepository.SearchByUsernameIgnoreCase(ctx, query)
	if users == nil {
		return nil, errors.NewStatusError(http.StatusInternalServerError, "Ошибка поиска")
	}

	return mapper.MapUserListToProxy(*users), nil
}
