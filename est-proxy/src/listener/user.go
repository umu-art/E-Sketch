package listener

import (
	"est-proxy/src/config"
	"est-proxy/src/models"
	"est-proxy/src/service"
	proxymodels "est_proxy_go/models"
	"fmt"

	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserListener struct {
	userRepository *service.UserRepository
}

func NewUserListener(userRepository *service.UserRepository) *UserListener {
	return &UserListener{userRepository: userRepository}
}

func (u UserListener) getSession(ctx echo.Context) (*uuid.UUID, error) {
	session := service.GetAndParseUserJWTCookie(ctx)
	if session == nil {
		return nil, ctx.String(http.StatusUnauthorized, fmt.Sprintf("Отсутствует или некорректная сессия"))
	}

	if session.ExpirationTime.After(time.Now()) {
		return nil, ctx.String(http.StatusUnauthorized, "Срок сессии истёк")
	}

	//TODO exists by id
	if user := u.userRepository.GetUserByID(&session.UserID); user == nil {
		return nil, ctx.String(http.StatusUnauthorized, "Пользователь не найден")
	}

	return &session.UserID, nil
}

func (u UserListener) CheckSession(ctx echo.Context) error {
	_, sessionStatus := u.getSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}
	return ctx.String(http.StatusOK, "Сессия действительна")
}

func (u UserListener) GetSelf(ctx echo.Context) error {
	userId, sessionStatus := u.getSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}

	user := u.userRepository.GetUserByID(userId)
	if user == nil {
		return ctx.String(http.StatusBadRequest, "Пользователь не найден")
	}

	return ctx.JSON(http.StatusOK, mapToProxyDto(user))
}

func (u UserListener) GetUserById(ctx echo.Context) error {
	_, sessionStatus := u.getSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}

	userIdStr := ctx.QueryParam("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорретный запрос: %v", err))
	}

	user := u.userRepository.GetUserByID(&userId)
	if user == nil {
		return ctx.String(http.StatusBadRequest, "Пользователь не найден")
	}

	return ctx.JSON(http.StatusOK, mapToProxyDto(user))
}

func (u UserListener) Login(ctx echo.Context) error {
	var authDto proxymodels.AuthDto
	if err := ctx.Bind(&authDto); err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорреткный запрос: %v", err))
	}

	user := u.userRepository.GetUserByEmail(authDto.Email)
	if user == nil {
		return ctx.String(http.StatusBadRequest, "Отсутвует или некорректный адрес почты или пароль")
	}

	if authDto.PasswordHash != user.PasswordHash {
		return ctx.String(http.StatusBadRequest, "Отсутвует или некорректный адрес почты или пароль")
	}

	token := service.GenerateUserJWTString(&user.ID)
	if token == nil {
		return ctx.String(http.StatusUnauthorized, "Не получилось начать сессию")
	}

	cookie := new(http.Cookie)
	cookie.Name = config.JWT_COOKIE_NAME
	cookie.Value = *token
	ctx.SetCookie(cookie)

	return ctx.String(http.StatusOK, "Вход в аккаунт выполнен успешно")
}

func (u UserListener) Register(ctx echo.Context) error {
	var regDto proxymodels.RegisterDto
	if err := ctx.Bind(&regDto); err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорреткный запрос: %v", err))
	}

	exists := u.userRepository.UserExistsByUsernameOrEmail(regDto.Username, regDto.Email)
	if exists == false {
		return ctx.String(http.StatusConflict, "Занято имя пользователя или адрес электронной почты")
	}

	userId := u.userRepository.Create(regDto.Username, regDto.Email, regDto.PasswordHash)
	if userId == nil {
		return ctx.String(http.StatusInternalServerError, "Не получилось создать аккаунт")
	}

	token := service.GenerateUserJWTString(userId)
	if token == nil {
		return ctx.String(http.StatusUnauthorized, "Не получилось начать сессию")
	}

	cookie := new(http.Cookie)
	cookie.Name = config.JWT_COOKIE_NAME
	cookie.Value = *token
	ctx.SetCookie(cookie)

	return ctx.String(http.StatusOK, "Аккаунт успешно зарегистрирован")
}

func (u UserListener) Search(ctx echo.Context) error {
	username := ctx.QueryParam("username")
	users := u.userRepository.SearchByUsernameIgnoreCase(ctx.Request().Context(), username)
	if users == nil {
		return ctx.String(http.StatusInternalServerError, "Ошибка поиска")
	}

	return ctx.JSON(http.StatusOK, mapListToProxyDto(*users))
}

func mapToProxyDto(user *models.User) proxymodels.UserDto {
	return proxymodels.UserDto{
		Id:       user.ID.String(),
		Username: user.Username,
		Avatar:   user.Avatar,
	}
}

func mapListToProxyDto(users []models.User) []proxymodels.UserDto {
	var dtos []proxymodels.UserDto
	for _, user := range users {
		dtos = append(dtos, mapToProxyDto(&user))
	}
	return dtos
}
