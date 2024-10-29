package service

import (
	"est-proxy/src/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const kSecretKey string = "real-secret-for-real-men"
const kDefaultJWTDuration time.Duration = time.Hour * 2

func ParseUserJWTtoken(token jwt.Token) (models.UserToken, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userUUIDStr, ok := claims["userUUID"].(string)
		if !ok {
			return models.UserToken{}, fmt.Errorf("missing or invalid 'userUUID' claim")
		}
		userUUID, err := uuid.Parse(userUUIDStr)
		if err != nil {
			return models.UserToken{}, fmt.Errorf("invalid 'userUUID' format: %w", err)
		}

		expirationTimeFloat, ok := claims["expirationTime"].(float64)
		if !ok {
			return models.UserToken{}, fmt.Errorf("missing or invalid 'expirationTime' claim")
		}
		expirationTime := time.Unix(int64(expirationTimeFloat), 0)

		return models.UserToken{
			UserUUID:       userUUID,
			ExpirationTime: expirationTime,
		}, nil
	} else {
		return models.UserToken{}, fmt.Errorf("ivanild JWT claiming")
	}
}

func GetUserJWTtoken(ctx echo.Context) (*jwt.Token, error) {
	user, ok := ctx.Get("user").(*jwt.Token)
	if !ok {
		return nil, fmt.Errorf("missing or invalid user in context")
	}
	return user, nil
}

func GetAndParseUserJWT(ctx echo.Context) (models.UserToken, error) {
	jwtToken, jwt_err := GetUserJWTtoken(ctx)
	if jwt_err != nil {
		return models.UserToken{}, jwt_err
	}
	return ParseUserJWTtoken(*jwtToken)
}

func GenerateUserJWTtoken(userUUID uuid.UUID) (jwt.Token, error) {
	expirationTime := time.Now().Add(kDefaultJWTDuration)

	claims := jwt.MapClaims{
        "userUUID":       userUUID.String(),
        "expirationTime": expirationTime.Unix(),
    }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return *token, nil
}

func GenerateUserJWTstring(userUUID uuid.UUID) (string, error) {
	token, err := GenerateUserJWTtoken(userUUID)
	if err != nil {
		return "", err
	}
	tokenString, err := token.SignedString([]byte(kSecretKey))
    if err != nil {
        return "", err
    }
	return tokenString, nil
}