package service

import (
	"est-proxy/src/config"
	"est-proxy/src/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
)

func SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if isExludedSessionPath(ctx.Path()) {
			return next(ctx)
		}

		session := utils.GetAndParseUserJWTCookie(ctx)

		if session == nil || session.ExpirationTime.Before(time.Now()) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
		}

		ctx.Set("sessionUserId", session.UserID.String())

		return next(ctx)
	}
}

func GetSessionUserId(ctx echo.Context) *uuid.UUID {
	userIdStr := ctx.Get("sessionUserId")
	if userIdStr == nil {
		return nil
	}

	userId, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		return nil
	}

	return &userId
}

func isExludedSessionPath(path string) bool {
	suffixIndex := strings.LastIndex(path, "/")
	if suffixIndex == -1 {
		return false
	}
	suffix := path[suffixIndex:]
	for _, v := range config.SESSION_CHECK_EXCLUDED_PATH_SUFFIXES {
		if suffix == v {
			return true
		}
	}
	return false
}
