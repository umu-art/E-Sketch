package listener

import (
	"est-proxy/src/config"
	"est-proxy/src/service"
	proxymodels "est_proxy_go/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserListener struct {
	userService *service.UserService
}

func NewUserListener(userService *service.UserService) *UserListener {
	return &UserListener{userService: userService}
}

func (u UserListener) CheckSession(ctx echo.Context) error {
	_, err := u.userService.GetSession(ctx)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, err.Error())
	}

	return ctx.String(http.StatusOK, "Сессия действительна")
}

func (u UserListener) GetSelf(ctx echo.Context) error {
	userId, err := u.userService.GetSession(ctx)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, err.Error())
	}

	user, authStatus := u.userService.GetUserById(ctx, userId)
	if authStatus != nil {
		return authStatus
	}

	return ctx.JSON(http.StatusOK, *user)
}

func (u UserListener) GetUserById(ctx echo.Context) error {
	_, err := u.userService.GetSession(ctx)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, err.Error())
	}

	userIdStr := ctx.QueryParam("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорретный запрос: %v", err))
	}

	user, authStatus := u.userService.GetUserById(ctx, &userId)
	if authStatus != nil {
		return authStatus
	}

	return ctx.JSON(http.StatusOK, *user)
}

func (u UserListener) Login(ctx echo.Context) error {
	var authDto proxymodels.AuthDto
	if err := ctx.Bind(&authDto); err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорреткный запрос: %v", err))
	}

	token, authStatus := u.userService.Login(ctx, &authDto)
	if authStatus != nil {
		return authStatus
	}

	cookie := new(http.Cookie)
	cookie.Name = config.JWT_COOKIE_NAME
	cookie.Value = *token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	ctx.SetCookie(cookie)

	return ctx.String(http.StatusOK, "Вход в аккаунт выполнен успешно")
}

func (u UserListener) Register(ctx echo.Context) error {
	var regDto proxymodels.RegisterDto
	if err := ctx.Bind(&regDto); err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорреткный запрос: %v", err))
	}

	token, authStatus := u.userService.Register(ctx, &regDto)
	if authStatus != nil {
		return authStatus
	}

	cookie := new(http.Cookie)
	cookie.Name = config.JWT_COOKIE_NAME
	cookie.Value = *token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	ctx.SetCookie(cookie)

	return ctx.String(http.StatusOK, "Аккаунт успешно зарегистрирован")
}

func (u UserListener) Search(ctx echo.Context) error {
	username := ctx.QueryParam("username")
	users, authStatus := u.userService.Search(ctx, username)
	if authStatus != nil {
		return authStatus
	}

	return ctx.JSON(http.StatusOK, *users)
}
