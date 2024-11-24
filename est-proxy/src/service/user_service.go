package service

import (
	"context"
	"est-proxy/src/models"
	"est-proxy/src/models/errors"
	proxymodels "est_proxy_go/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type UserService struct {
	userRepository *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (u UserService) GetSession(ctx echo.Context) (*uuid.UUID, *errors.StatusError) {
	session := GetAndParseUserJWTCookie(ctx)

	if session == nil || session.ExpirationTime.Before(time.Now()) {
		return nil, errors.NewStatusError(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	return &session.UserID, nil
}

func (u UserService) GetUserById(ctx context.Context, userId *uuid.UUID) (*proxymodels.UserDto, *errors.StatusError) {
	user := u.userRepository.GetUserByID(ctx, userId)
	if user == nil {
		return nil, errors.NewStatusError(http.StatusBadRequest, "Пользователь не найден")
	}

	return mapToProxyDto(user.Public()), nil
}

func (u UserService) Login(ctx context.Context, authDto *proxymodels.AuthDto) (*string, *errors.StatusError) {
	user := u.userRepository.GetUserByEmail(ctx, authDto.Email)
	if user == nil {
		return nil, errors.NewStatusError(http.StatusBadRequest, "Отсутвует или некорректный адрес почты или пароль")
	}

	if authDto.PasswordHash != user.PasswordHash {
		return nil, errors.NewStatusError(http.StatusBadRequest, "Отсутвует или некорректный адрес почты или пароль")
	}

	token := GenerateUserJWTString(&user.ID)
	if token == nil {
		return nil, errors.NewStatusError(http.StatusUnauthorized, "Не получилось начать сессию")
	}

	return token, nil
}

func (u UserService) Register(ctx context.Context, registerDto *proxymodels.RegisterDto) (*string, *errors.StatusError) {
	exists := u.userRepository.UserExistsByUsernameOrEmail(ctx, registerDto.Username, registerDto.Email)
	if exists {
		return nil, errors.NewStatusError(http.StatusConflict, "Занято имя пользователя или адрес электронной почты")
	}

	userId := u.userRepository.Create(ctx, registerDto.Username, registerDto.Email, registerDto.PasswordHash)
	if userId == nil {
		return nil, errors.NewStatusError(http.StatusInternalServerError, "Не получилось создать аккаунт")
	}

	token := GenerateUserJWTString(userId)
	if token == nil {
		return nil, errors.NewStatusError(http.StatusUnauthorized, "Не получилось начать сессию")
	}

	return token, nil
}

func (u UserService) Search(ctx context.Context, query string) (*[]proxymodels.UserDto, *errors.StatusError) {
	users := u.userRepository.SearchByUsernameIgnoreCase(ctx, query)
	if users == nil {
		return nil, errors.NewStatusError(http.StatusInternalServerError, "Ошибка поиска")
	}

	return mapListToProxyDto(*users), nil
}

func mapToProxyDto(user *models.PublicUser) *proxymodels.UserDto {
	return &proxymodels.UserDto{
		Id:       user.ID.String(),
		Username: user.Username,
		Avatar:   user.Avatar,
	}
}

func mapListToProxyDto(users []models.PublicUser) *[]proxymodels.UserDto {
	var dtos []proxymodels.UserDto
	for _, user := range users {
		dtos = append(dtos, *mapToProxyDto(&user))
	}
	return &dtos
}
