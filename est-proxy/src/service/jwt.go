package service

import (
	"est-proxy/src/config"
	"est-proxy/src/models"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ParseUserJWT(token *jwt.Token) *models.ParsedJWT {
	if token == nil {
		return nil
	}
	if token.Method.Alg() != config.JWT_SIGNING_METHOD {
		log.Printf("Unsupported signing method: %v", token.Header["alg"])
		return nil
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userIDStr, ok := claims["userID"].(string)
		if !ok {
			log.Printf("Failed to parse user id from claims: %v", claims)
			return nil
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			log.Printf("Failed to parse user id from claims: %v", claims)
			return nil
		}
		expirationTimeFloat, ok := claims["exp"].(float64)
		if !ok {
			log.Printf("Failed to parse expiration time from claims: %v", claims)
			return nil
		}
		expirationTimeUnix := int64(expirationTimeFloat)
		expirationTime := time.Unix(expirationTimeUnix, 0)

		return &models.ParsedJWT{
			UserID:         userID,
			ExpirationTime: expirationTime,
		}
	}
	return nil
}

func GetUserJWTCookie(ctx echo.Context) *jwt.Token {
	tokenCookie, err := ctx.Cookie(config.JWT_COOKIE_NAME)
	if err != nil {
		log.Printf("Failed to get cookie %v", err)
		return nil
	}

	token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != config.JWT_SIGNING_METHOD {
			log.Printf("Invalid signing method: %s", token.Header["alg"])
			return nil, fmt.Errorf("неподдерживаемый алгоритм подписи")
		}
		return []byte(config.JWT_SECRET), nil
	})

	if err != nil {
		log.Printf("Failed to parse cookie %v", err)
		return nil
	}
	return token
}

func GetAndParseUserJWTCookie(ctx echo.Context) *models.ParsedJWT {
	jwtToken := GetUserJWTCookie(ctx)
	if jwtToken == nil {
		return nil
	}
	return ParseUserJWT(jwtToken)
}

func GenerateUserJWT(userID *uuid.UUID) (*jwt.Token, error) {
	expirationTime := time.Now().Add(config.JWT_DURATION_TIME)

	claims := jwt.MapClaims{
		"userID": userID.String(),
		"exp":    expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token, nil
}

func GenerateUserJWTString(userID *uuid.UUID) *string {
	token, err := GenerateUserJWT(userID)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return nil
	}
	tokenString, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		log.Printf("Failed to sign token: %v", err)
		return nil
	}
	return &tokenString
}
