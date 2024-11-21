package service

import (
	"est-proxy/src/models"
	proxymodels "est_proxy_go/models"
	"fmt"
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

func (u UserService) GetSession(ctx echo.Context) (*uuid.UUID, error) {
	session := GetAndParseUserJWTCookie(ctx)

	if session == nil || session.ExpirationTime.After(time.Now()) {
		return nil, ctx.String(http.StatusUnauthorized, fmt.Sprintf("Отсутствует или некорректная сессия"))
	}

	return &session.UserID, nil
}

func (u UserService) GetUserById(ctx echo.Context, userId *uuid.UUID) (*proxymodels.UserDto, error) {
	user := u.userRepository.GetUserByID(userId)
	if user == nil {
		return nil, ctx.String(http.StatusBadRequest, "Пользователь не найден")
	}

	return mapToProxyDto(user.Public()), nil
}

func (u UserService) Login(ctx echo.Context, authDto *proxymodels.AuthDto) (*string, error) {
	user := u.userRepository.GetUserByEmail(authDto.Email)
	if user == nil {
		return nil, ctx.String(http.StatusBadRequest, "Отсутвует или некорректный адрес почты или пароль")
	}

	if authDto.PasswordHash != user.PasswordHash {
		return nil, ctx.String(http.StatusBadRequest, "Отсутвует или некорректный адрес почты или пароль")
	}

	token := GenerateUserJWTString(&user.ID)
	if token == nil {
		return nil, ctx.String(http.StatusUnauthorized, "Не получилось начать сессию")
	}

	return token, nil
}

func (u UserService) Register(ctx echo.Context, registerDto *proxymodels.RegisterDto) (*string, error) {
	exists := u.userRepository.UserExistsByUsernameOrEmail(registerDto.Username, registerDto.Email)
	if exists {
		return nil, ctx.String(http.StatusConflict, "Занято имя пользователя или адрес электронной почты")
	}

	userId := u.userRepository.Create(registerDto.Username, registerDto.Email, registerDto.PasswordHash)
	if userId == nil {
		return nil, ctx.String(http.StatusInternalServerError, "Не получилось создать аккаунт")
	}

	token := GenerateUserJWTString(userId)
	if token == nil {
		return nil, ctx.String(http.StatusUnauthorized, "Не получилось начать сессию")
	}

	return token, nil
}

func (u UserService) Search(ctx echo.Context, query string) (*[]proxymodels.UserDto, error) {
	users := u.userRepository.SearchByUsernameIgnoreCase(ctx.Request().Context(), query)
	if users == nil {
		return nil, ctx.String(http.StatusInternalServerError, "Ошибка поиска")
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
