package utils

import (
	"est-proxy/src/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseUserJWT(t *testing.T) {
	t.Run("Нулевой токен", func(t *testing.T) {
		result := ParseUserJWT(nil)
		assert.Nil(t, result)
	})

	t.Run("Некорректный алгоритм шифрования", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS512)
		result := ParseUserJWT(token)
		assert.Nil(t, result)
	})

	t.Run("Отсутствует UserId", func(t *testing.T) {
		claims := jwt.MapClaims{
			"exp": float64(time.Now().Add(time.Hour).Unix()),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		result := ParseUserJWT(token)

		assert.Nil(t, result)
	})

	t.Run("Некорректный UserID", func(t *testing.T) {
		claims := jwt.MapClaims{
			"userID": 123,
			"exp":    float64(time.Now().Add(time.Hour).Unix()),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		result := ParseUserJWT(token)
		assert.Nil(t, result)
	})

	t.Run("Корректный токен", func(t *testing.T) {
		userID := uuid.New()
		exp := time.Now().Add(time.Hour).Unix()

		claims := jwt.MapClaims{
			"userID": userID.String(),
			"exp":    float64(exp),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		result := ParseUserJWT(token)
		require.NotNil(t, result)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, exp, result.ExpirationTime.Unix())
	})
}

func TestGetUserJWTCookie(t *testing.T) {
	e := echo.New()

	t.Run("Отсутствует Cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		result := GetUserJWTCookie(c)
		assert.Nil(t, result)
	})

	t.Run("Некорректный токен", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  config.JWT_COOKIE_NAME,
			Value: "invalid.token.value",
		})
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		result := GetUserJWTCookie(c)
		assert.Nil(t, result)
	})

	t.Run("Корректный токен", func(t *testing.T) {
		userID := uuid.New()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID": userID.String(),
			"exp":    time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(config.JWT_SECRET))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  config.JWT_COOKIE_NAME,
			Value: tokenString,
		})
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		result := GetUserJWTCookie(c)
		require.NotNil(t, result)
		assert.True(t, result.Valid)
	})
}

func TestGetAndParseUserJWTCookie(t *testing.T) {
	t.Run("Отсутствует Cookie", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		result := GetAndParseUserJWTCookie(c)
		assert.Nil(t, result)
	})

	t.Run("Корректная Cookie", func(t *testing.T) {
		userID := uuid.New()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID": userID.String(),
			"exp":    time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(config.JWT_SECRET))

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  config.JWT_COOKIE_NAME,
			Value: tokenString,
		})
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		result := GetAndParseUserJWTCookie(c)
		require.NotNil(t, result)
		assert.Equal(t, userID, result.UserID)
	})
}

func TestGenerateUserJWT(t *testing.T) {
	userID := uuid.New()
	token, err := GenerateUserJWT(&userID)
	require.NoError(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, userID.String(), claims["userID"])
	assert.IsType(t, int64(0), claims["exp"])
	assert.Equal(t, jwt.SigningMethodHS256, token.Method)
}

func TestGenerateUserJWTString(t *testing.T) {
	t.Run("Некорректный ключ шифрования", func(t *testing.T) {
		originalSecret := config.JWT_SECRET
		defer func() { config.JWT_SECRET = originalSecret }()

		config.JWT_SECRET = "different_secret"
		userID := uuid.New()
		tokenStr := GenerateUserJWTString(&userID)
		require.NotNil(t, tokenStr)

		config.JWT_SECRET = originalSecret
		_, err := jwt.Parse(*tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.JWT_SECRET), nil
		})
		assert.Error(t, err)
	})

	t.Run("Корректный ключ", func(t *testing.T) {
		userID := uuid.New()
		tokenStr := GenerateUserJWTString(&userID)
		require.NotNil(t, tokenStr)
		assert.NotEmpty(t, *tokenStr)

		token, err := jwt.Parse(*tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.JWT_SECRET), nil
		})
		require.NoError(t, err)
		assert.True(t, token.Valid)
	})
}
