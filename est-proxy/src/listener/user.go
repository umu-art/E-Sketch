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
	authService *service.AuthService
}

func NewUserListener(authService *service.AuthService) *UserListener {
	return &UserListener{authService: authService}
}

func (u UserListener) CheckSession(ctx echo.Context) error {
	_, sessionStatus := u.authService.GetSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}
	return ctx.String(http.StatusOK, "Сессия действительна")
}

func (u UserListener) GetSelf(ctx echo.Context) error {
	userId, sessionStatus := u.authService.GetSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}

	user, authStatus := u.authService.GetUserById(ctx, userId)
	if authStatus != nil {
		return authStatus
	}

	return ctx.JSON(http.StatusOK, *user)
}

func (u UserListener) GetUserById(ctx echo.Context) error {
	_, sessionStatus := u.authService.GetSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}

	userIdStr := ctx.QueryParam("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорретный запрос: %v", err))
	}

	user, authStatus := u.authService.GetUserById(ctx, &userId)
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

	token, authStatus := u.authService.Login(ctx, &authDto)
	if authStatus != nil {
		return authStatus
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

	token, authStatus := u.authService.Register(ctx, &regDto)
	if authStatus != nil {
		return authStatus
	}

	cookie := new(http.Cookie)
	cookie.Name = config.JWT_COOKIE_NAME
	cookie.Value = *token
	ctx.SetCookie(cookie)

	return ctx.String(http.StatusOK, "Аккаунт успешно зарегистрирован")
}

func (u UserListener) Search(ctx echo.Context) error {
	username := ctx.QueryParam("username")
	users, authStatus := u.authService.Search(ctx, username)
	if authStatus != nil {
		return authStatus
	}

	return ctx.JSON(http.StatusOK, *users)
}
