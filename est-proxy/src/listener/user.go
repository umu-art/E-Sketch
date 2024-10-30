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

func (u UserListener) getSessionStatus(ctx echo.Context) int {
	session, err := service.GetAndParseUserJWT(ctx)
	if err != nil {
		return http.StatusBadRequest
	}

	if session.ExpirationTime.After(time.Now()) {
		return http.StatusUnauthorized
	}

	if _, err := u.userRepository.GetUserByID(&session.UserID); err != nil {
		return http.StatusUnauthorized
	}

	return http.StatusOK
}

func (u UserListener) CheckSession(ctx echo.Context) error {
	return ctx.NoContent(u.getSessionStatus(ctx))
}

func (u UserListener) GetSelf(ctx echo.Context) error {
	// TODO нужно ли вообще проверять сессию? 
	
	// sessionStatus := u.getSessionStatus(ctx)
	// if sessionStatus != http.StatusOK {
	// 	return ctx.NoContent(http.StatusUnauthorized)
	// }

	session, err := service.GetAndParseUserJWT(ctx)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	
	user, err := u.userRepository.GetUserByID(&session.UserID)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, mapToProxyDTO(user))
}

func (u UserListener) GetUserById(ctx echo.Context) error {
	userIdStr := ctx.QueryParam("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	user, err := u.userRepository.GetUserByID(&userId)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, mapToProxyDTO(user))
}

func (u UserListener) Login(ctx echo.Context) error {
	var authDto proxy_models.AuthDto
	if err := ctx.Bind(&authDto); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	userId, err := u.userRepository.GetIDByUserInfo(&authDto.Username, nil)
	if err != nil {
		// TODO: нужно отправлять ошибку или просто статус ОК без JWT? 
		return ctx.NoContent(http.StatusUnauthorized)
	}
	
	user, err := u.userRepository.GetUserByID(userId)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	if authDto.PasswordHash != user.PasswordHash {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	token, err := service.GenerateUserJWTstring(*userId)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	ctx.Response().Header().Set("Authorization", "Bearer " + token)
	return ctx.NoContent(http.StatusOK)
}

func (u UserListener) Register(ctx echo.Context) error {
	var regDto proxy_models.RegisterDto
	if err := ctx.Bind(&regDto); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	userId, _ := u.userRepository.GetIDByUserInfo(&regDto.Username, &regDto.Email)
	if userId != nil {
		return ctx.NoContent(http.StatusConflict)
	}

	err := u.userRepository.Create(regDto.Username, regDto.Email, regDto.PasswordHash)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	userId, err = u.userRepository.GetIDByUserInfo(&regDto.Username, nil)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	token, err := service.GenerateUserJWTstring(*userId)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	ctx.Response().Header().Set("Authorization", "Bearer " + token)
	return ctx.NoContent(http.StatusOK)
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