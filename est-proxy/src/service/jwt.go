package service

import (
	"est-proxy/src/config"
	"est-proxy/src/models"

	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ParseUserJWTtoken(token jwt.Token) (*models.ParsedJWT, error) {
	if token.Method.Alg() != config.JWT_SIGNING_METHOD {
		return nil, fmt.Errorf("Invalid JWT signing method: %s, must be %s", token.Method.Alg(), config.JWT_SIGNING_METHOD)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userIDStr, ok := claims["userID"].(string)
		if !ok {
			return nil, fmt.Errorf("Missing or invalid 'userID' claim")
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("Invalid 'userID' format: %w", err)
		}

		expirationTimeUnix, ok := claims["expirationTime"].(int64)
		if !ok {
			return nil, fmt.Errorf("Missing or invalid 'expirationTime' claim")
		}
		expirationTime := time.Unix(expirationTimeUnix, 0)

		return &models.ParsedJWT{
			UserID:         userID,
			ExpirationTime: expirationTime,
		}, nil
	}
	return nil, fmt.Errorf("Failed to claim JWT")
}

func GetUserJWTtoken(ctx echo.Context) (*jwt.Token, error) {
	tokenCookie, err := ctx.Cookie("jwt_token")
	if err != nil {
		return nil, fmt.Errorf("Ошибка при попытке получить кукисы: %w", err)
	}
	
	token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != config.JWT_SIGNING_METHOD {
			return nil, fmt.Errorf("Неподдерживаемый алгоритм подписи")
		}
		return config.JWT_SECRET, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("Отсутствует или некорретный jwt токен: %w", err)
	}
	return token, nil
}

func GetAndParseUserJWT(ctx echo.Context) (*models.ParsedJWT, error) {
	jwtToken, err := GetUserJWTtoken(ctx)
	if err != nil {
		return nil, err
	}
	return ParseUserJWTtoken(*jwtToken)
}

func GenerateUserJWTtoken(userID uuid.UUID) (*jwt.Token, error) {
	expirationTime := time.Now().Add(config.JWT_DURATION_TIME)

	claims := jwt.MapClaims{
		"userID":       userID.String(),
		"expirationTime": expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token, nil
}

func GenerateUserJWTstring(userID uuid.UUID) (string, error) {
	token, err := GenerateUserJWTtoken(userID)
	if err != nil {
		return "", fmt.Errorf("Failed to generate token: %w", err)
	}
	tokenString, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		return "", fmt.Errorf("Failed to convert token to string: %w", err)
	}
	return tokenString, nil
}