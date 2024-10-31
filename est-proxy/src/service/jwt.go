package service

import (
	"est-proxy/src/models"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ParseUserJWTtoken(token jwt.Token) (*models.ParsedJWT, error) {
	signingMethod := os.Getenv("JWT_SIGNING_METHOD")

	if token.Method.Alg() != signingMethod {
		return nil, fmt.Errorf("Invalid JWT signing method: %s, must be %s", token.Method.Alg(), signingMethod)
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
	user, ok := ctx.Get("user").(*jwt.Token)
	if !ok {
		return nil, fmt.Errorf("Missing or invalid user in context")
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

func GenerateUserJWTtoken(userID uuid.UUID) (*jwt.Token, error) {
	durationStr := os.Getenv("JWT_DURATION_TIME")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to get jwt duration time: %w", err)
	}

	expirationTime := time.Now().Add(duration)

	claims := jwt.MapClaims{
		"userID":       userID.String(),
		"expirationTime": expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token, nil
}

func GenerateUserJWTstring(userID uuid.UUID) (string, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")

	token, err := GenerateUserJWTtoken(userID)
	if err != nil {
		return "", fmt.Errorf("Failed to generate token: %w", err)
	}
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("Failed to convert token to string: %w", err)
	}
	return tokenString, nil
}