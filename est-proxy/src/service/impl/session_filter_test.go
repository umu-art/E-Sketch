package impl

import (
	"est-proxy/src/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"est-proxy/src/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSessionMiddleware(t *testing.T) {
	e := echo.New()

	t.Run("excluded path skips session check", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/login", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/login")

		originalExcluded := config.SESSION_CHECK_EXCLUDED_PATH_SUFFIXES
		config.SESSION_CHECK_EXCLUDED_PATH_SUFFIXES = []string{"/login"}
		defer func() { config.SESSION_CHECK_EXCLUDED_PATH_SUFFIXES = originalExcluded }()

		called := false
		mockNext := func(ctx echo.Context) error {
			called = true
			return nil
		}

		handler := SessionMiddleware(mockNext)
		err := handler(ctx)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("valid session continues flow", func(t *testing.T) {
		userID := uuid.New()
		token := utils.GenerateUserJWTString(&userID)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  config.JWT_COOKIE_NAME,
			Value: *token,
		})
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/protected")

		called := false
		mockNext := func(ctx echo.Context) error {
			called = true
			return nil
		}

		handler := SessionMiddleware(mockNext)
		err := handler(ctx)

		assert.NoError(t, err)
		assert.True(t, called)
		assert.Equal(t, userID.String(), ctx.Get("sessionUserId"))
	})

	t.Run("expired session returns unauthorized", func(t *testing.T) {
		config.JWT_DURATION_TIME = -1 * time.Hour
		defer func() { config.JWT_DURATION_TIME = 1 * time.Hour }()

		userID := uuid.New()
		token := utils.GenerateUserJWTString(&userID)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  config.JWT_COOKIE_NAME,
			Value: *token,
		})
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/protected")

		handler := SessionMiddleware(func(ctx echo.Context) error {
			return nil
		})

		err := handler(ctx)
		assert.Error(t, err)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, he.Code)
	})

	t.Run("missing cookie returns unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/protected")

		handler := SessionMiddleware(func(ctx echo.Context) error {
			return nil
		})

		err := handler(ctx)
		assert.Error(t, err)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, he.Code)
	})
}

func TestGetSessionUserId(t *testing.T) {
	e := echo.New()

	t.Run("no user id in context", func(t *testing.T) {
		ctx := e.NewContext(nil, nil)
		result := GetSessionUserId(ctx)
		assert.Nil(t, result)
	})

	t.Run("invalid uuid format", func(t *testing.T) {
		ctx := e.NewContext(nil, nil)
		ctx.Set("sessionUserId", "invalid-uuid")
		result := GetSessionUserId(ctx)
		assert.Nil(t, result)
	})

	t.Run("valid uuid in context", func(t *testing.T) {
		ctx := e.NewContext(nil, nil)
		expected := uuid.New()
		ctx.Set("sessionUserId", expected.String())
		result := GetSessionUserId(ctx)
		assert.Equal(t, expected.String(), result.String())
	})
}

func TestIsExludedSessionPath(t *testing.T) {
	t.Run("excluded suffix", func(t *testing.T) {
		originalExcluded := config.SESSION_CHECK_EXCLUDED_PATH_SUFFIXES
		config.SESSION_CHECK_EXCLUDED_PATH_SUFFIXES = []string{"/public"}
		defer func() { config.SESSION_CHECK_EXCLUDED_PATH_SUFFIXES = originalExcluded }()

		assert.True(t, isExludedSessionPath("/api/v1/public"))
	})

	t.Run("non-excluded suffix", func(t *testing.T) {
		originalExcluded := config.SESSION_CHECK_EXCLUDED_PATH_SUFFIXES
		config.SESSION_CHECK_EXCLUDED_PATH_SUFFIXES = []string{"/public"}
		defer func() { config.SESSION_CHECK_EXCLUDED_PATH_SUFFIXES = originalExcluded }()

		assert.False(t, isExludedSessionPath("/api/v1/private"))
	})

	t.Run("root path excluded", func(t *testing.T) {
		originalExcluded := config.SESSION_CHECK_EXCLUDED_PATH_SUFFIXES
		config.SESSION_CHECK_EXCLUDED_PATH_SUFFIXES = []string{"/"}
		defer func() { config.SESSION_CHECK_EXCLUDED_PATH_SUFFIXES = originalExcluded }()

		assert.True(t, isExludedSessionPath("/"))
	})
}
