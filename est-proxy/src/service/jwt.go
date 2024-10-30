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
const kSigningMethog string = "HS256"
const kDefaultJWTDuration time.Duration = time.Hour * 2

func ParseUserJWTtoken(token jwt.Token) (*models.ParsedJWT, error) {
	if token.Method.Alg() != kSigningMethog {
		return nil, fmt.Errorf("invalid JWT signing method: %w, must be %w", token.Method.Alg(), kSigningMethog)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userIDStr, ok := claims["userID"].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid 'userUUID' claim")
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid 'userID' format: %w", err)
		}

		expirationTimeUnix, ok := claims["expirationTime"].(int64)
		if !ok {
			return nil, fmt.Errorf("missing or invalid 'expirationTime' claim")
		}
		expirationTime := time.Unix(expirationTimeUnix, 0)

		return &models.ParsedJWT{
			UserID:         userID,
			ExpirationTime: expirationTime,
		}, nil
	}
	return nil, fmt.Errorf("failed to claim JWT")
}

func GetUserJWTtoken(ctx echo.Context) (*jwt.Token, error) {
	user, ok := ctx.Get("user").(*jwt.Token)
	if !ok {
		return nil, fmt.Errorf("missing or invalid user in context")
	}
	return user, nil
}

func GetAndParseUserJWT(ctx echo.Context) (*models.ParsedJWT, error) {
	jwtToken, err := GetUserJWTtoken(ctx)
	if err != nil {
		return nil, err
	}
	return ParseUserJWTtoken(*jwtToken)
}

func GenerateUserJWTtoken(userID uuid.UUID) (jwt.Token, error) {
	expirationTime := time.Now().Add(kDefaultJWTDuration)

	claims := jwt.MapClaims{
		"userID":       userID.String(),
		"expirationTime": expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return *token, nil
}

func GenerateUserJWTstring(userID uuid.UUID) (string, error) {
	token, err := GenerateUserJWTtoken(userID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	tokenString, err := token.SignedString([]byte(kSecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to convert token to string: %w", err)
	}
	return tokenString, nil
}