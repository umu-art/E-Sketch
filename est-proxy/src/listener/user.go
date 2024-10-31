package listener

import (
	"est-proxy/src/models"
	"est-proxy/src/service"
	proxy_models "est_proxy_go/models"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserListener struct {
	userRepository *service.UserRepository
}

func NewUserListener(userRepository *service.UserRepository) *UserListener {
	return &UserListener{userRepository: userRepository}
}

func (u UserListener) getSessionStatus(ctx echo.Context) (int, string) {
	session, err := service.GetAndParseUserJWT(ctx)
	if err != nil {
		return http.StatusBadRequest, "Missing or invalid session token"
	}

	if session.ExpirationTime.After(time.Now()) {
		return http.StatusUnauthorized, "Session expired"
	}

	if _, err := u.userRepository.GetUserByID(&session.UserID); err != nil {
		return http.StatusUnauthorized, "User entity not found"
	}

	return http.StatusOK, "Session if valid"
}

func (u UserListener) CheckSession(ctx echo.Context) error {
	statusCode, statusMessage := u.getSessionStatus(ctx)
	return ctx.String(statusCode, statusMessage)
}

func (u UserListener) GetSelf(ctx echo.Context) error {
	session, err := service.GetAndParseUserJWT(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Missing or invalid session token")
	}
	
	user, err := u.userRepository.GetUserByID(&session.UserID)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "User entity not found")
	}

	return ctx.JSON(http.StatusOK, mapToProxyDTO(user))
}

func (u UserListener) GetUserById(ctx echo.Context) error {
	userIdStr := ctx.QueryParam("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Missing or invalid User ID")
	}

	user, err := u.userRepository.GetUserByID(&userId)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "User entity not found")
	}

	return ctx.JSON(http.StatusOK, mapToProxyDTO(user))
}

func (u UserListener) Login(ctx echo.Context) error {
	var authDto proxy_models.AuthDto
	if err := ctx.Bind(&authDto); err != nil {
		return ctx.String(http.StatusBadRequest, "Missing or invalid email or password")
	}

	userId, err := u.userRepository.GetIDByEmail(authDto.Email)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "User entity not found")
	}
	
	user, err := u.userRepository.GetUserByID(userId)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "User entity not found")
	}

	if authDto.PasswordHash != user.PasswordHash {
		return ctx.String(http.StatusUnauthorized, "Incorrect password")
	}

	token, err := service.GenerateUserJWTstring(*userId)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to create session")
	}

	ctx.Response().Header().Set("Authorization", token)
	return ctx.String(http.StatusOK, "Account logged in successfully")
}

func (u UserListener) Register(ctx echo.Context) error {
	var regDto proxy_models.RegisterDto
	if err := ctx.Bind(&regDto); err != nil {
		return ctx.String(http.StatusBadRequest, "Missing or invalid username, email or password")
	}

	userId, _ := u.userRepository.GetIDByUsername(regDto.Username)
	if userId != nil {
		return ctx.String(http.StatusConflict, "Account with this username already exists")
	}

	userId, _ = u.userRepository.GetIDByEmail(regDto.Email)
	if userId != nil {
		return ctx.String(http.StatusConflict, "Account with this email already exists")
	}

	err := u.userRepository.Create(regDto.Username, regDto.Email, regDto.PasswordHash)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to create account")
	}

	userId, err = u.userRepository.GetIDByUsername(regDto.Username)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Created user is not found")
	}

	token, err := service.GenerateUserJWTstring(*userId)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to create session")
	}

	ctx.Response().Header().Set("Authorization", token)
	return ctx.String(http.StatusOK, "Account is registered successfully")
}

func (u UserListener) Search(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func mapToProxyDTO(user *models.User) proxy_models.UserDto {
	return proxy_models.UserDto{
		Id: user.ID.String(),
		Username: user.Username,
		Avatar: user.Avatar,
	}
}