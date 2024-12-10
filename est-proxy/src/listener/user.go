package listener

import (
	"est-proxy/src/config"
	"est-proxy/src/service"
	"est-proxy/src/service/impl"
	proxymodels "est_proxy_go/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserListener struct {
	userService service.UserService
}

func NewUserListener(userService service.UserService) *UserListener {
	return &UserListener{userService: userService}
}

func (u UserListener) CheckSession(ctx echo.Context) error {
	sessionUserId := impl.GetSessionUserId(ctx)
	if sessionUserId == nil {
		return ctx.String(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	return ctx.String(http.StatusOK, "Сессия действительна")
}

func (u UserListener) GetSelf(ctx echo.Context) error {
	sessionUserId := impl.GetSessionUserId(ctx)
	if sessionUserId == nil {
		return ctx.String(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	user, statusError := u.userService.GetUserById(ctx.Request().Context(), sessionUserId)
	if statusError != nil {
		return statusError.Send(ctx)
	}

	return ctx.JSON(http.StatusOK, *user)
}

func (u UserListener) GetUserById(ctx echo.Context) error {
	sessionUserId := impl.GetSessionUserId(ctx)
	if sessionUserId == nil {
		return ctx.String(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	userIdStr := ctx.Param("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорретный запрос: %v", err))
	}

	user, statusError := u.userService.GetUserById(ctx.Request().Context(), &userId)
	if statusError != nil {
		return statusError.Send(ctx)
	}

	return ctx.JSON(http.StatusOK, *user)
}

func (u UserListener) Login(ctx echo.Context) error {
	var authDto proxymodels.AuthDto
	if err := ctx.Bind(&authDto); err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорреткный запрос: %v", err))
	}

	token, statusError := u.userService.Login(ctx.Request().Context(), &authDto)
	if statusError != nil {
		return statusError.Send(ctx)
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

	token, statusError := u.userService.Register(ctx.Request().Context(), &regDto)
	if statusError != nil {
		return statusError.Send(ctx)
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
	sessionUserId := impl.GetSessionUserId(ctx)
	if sessionUserId == nil {
		return ctx.String(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	username := ctx.Param("username")
	users, statusError := u.userService.Search(ctx.Request().Context(), username)
	if statusError != nil {
		return statusError.Send(ctx)
	}

	return ctx.JSON(http.StatusOK, *users)
}
