package listener

import (
	"est-proxy/src/config"
	"est-proxy/src/models"
	"est-proxy/src/service"
	proxy_models "est_proxy_go/models"

	"errors"
	"time"
	
	"net/http"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserListener struct {
	userRepository *service.UserRepository
}

func NewUserListener(userRepository *service.UserRepository) *UserListener {
	return &UserListener{userRepository: userRepository}
}

func (u UserListener) getSession(ctx echo.Context) (*uuid.UUID, error) {
	session, err := service.GetAndParseUserJWT(ctx)
	if err != nil {
		return nil, ctx.String(http.StatusUnauthorized, "Отсутствует или некорректный токен")
	}

	if session.ExpirationTime.After(time.Now()) {
		return nil, ctx.String(http.StatusUnauthorized, "Срок сессии истёк")
	}

	if _, err := u.userRepository.GetUserByID(&session.UserID); err != nil {
		return nil, ctx.String(http.StatusUnauthorized, "Пользователь не найден")
	}

	return &session.UserID, nil
}

func (u UserListener) CheckSession(ctx echo.Context) error {
	_, sessionStatus := u.getSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}
	return ctx.String(http.StatusOK, "Сессия действительна")
}

func (u UserListener) GetSelf(ctx echo.Context) error {
	userId, sessionStatus := u.getSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}
	
	user, err := u.userRepository.GetUserByID(userId)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Пользователь не найден")
	}

	return ctx.JSON(http.StatusOK, mapToProxyDto(user))
}

func (u UserListener) GetUserById(ctx echo.Context) error {
	_, sessionStatus := u.getSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}
	
	userIdStr := ctx.QueryParam("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Отсутствует или некорректный ID")
	}

	user, err := u.userRepository.GetUserByID(&userId)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Пользователь не найден")
	}

	return ctx.JSON(http.StatusOK, mapToProxyDto(user))
}

func (u UserListener) Login(ctx echo.Context) error {
	var authDto proxy_models.AuthDto
	if err := ctx.Bind(&authDto); err != nil {
		return ctx.String(http.StatusBadRequest, "Отствует или некорректный адрес почты или пароль")
	}

	user, err := u.userRepository.GetUserByEmail(authDto.Email)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Отствует или некорректный адрес почты или пароль")
	}

	if authDto.PasswordHash != user.PasswordHash {
		return ctx.String(http.StatusBadRequest, "Отствует или некорректный адрес почты или пароль")
	}

	token, err := service.GenerateUserJWTstring(user.ID)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, "Не получилось начать сессию")
	}

	cookie := new(http.Cookie)
	cookie.Name = config.JWT_COOKIE_NAME
	cookie.Value = token
	ctx.SetCookie(cookie)

	return ctx.String(http.StatusOK, "Вход в аккаунт выполнен успешно")
}

func (u UserListener) Register(ctx echo.Context) error {
	var regDto proxy_models.RegisterDto
	if err := ctx.Bind(&regDto); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, sub_err := range ve {
				switch sub_err.Field() {
                case "Username":
                    return ctx.String(http.StatusBadRequest, "Неверное имя пользователя")
                case "Email":
                    return ctx.String(http.StatusBadRequest, "Неверный адрес электронной почты")
                case "PasswordHash":
                    return ctx.String(http.StatusBadRequest, "Неверный пароль")
                default:
                    return ctx.String(http.StatusBadRequest, "Неизвестная ошибка")
                }
			}
		}
		return ctx.String(http.StatusBadRequest, "Неизвестная ошибка")
	}

	exists, err := u.userRepository.UserExistsByUsernameOrEmail(regDto.Username, regDto.Email)
	if exists == false || err != nil {
		return ctx.String(http.StatusConflict, "Имя пользователя или адрес электронной почты уже занято")
	}

	err = u.userRepository.Create(regDto.Username, regDto.Email, regDto.PasswordHash)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Не получилось создать аккаунт")
	}

	userId, err := u.userRepository.GetIDByEmail(regDto.Email)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Пользователь не найден")
	}

	token, err := service.GenerateUserJWTstring(*userId)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, "Не получилось начать сессию")
	}

	cookie := new(http.Cookie)
	cookie.Name = config.JWT_COOKIE_NAME
	cookie.Value = token
	ctx.SetCookie(cookie)
	
	return ctx.String(http.StatusOK, "Аккаунт успешно зарегистрирован")
}

func (u UserListener) Search(ctx echo.Context) error {
	username := ctx.QueryParam("username")
	users, err := u.userRepository.SearchByUsernameIgnoreCase(ctx.Request().Context(), username)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Ошибка поиска")
	}

	return ctx.JSON(http.StatusOK, mapListToProxyDto(users))
}

func mapToProxyDto(user *models.User) proxy_models.UserDto {
	return proxy_models.UserDto{
		Id: user.ID.String(),
		Username: user.Username,
		Avatar: user.Avatar,
	}
}

func mapListToProxyDto(users []models.User) []proxy_models.UserDto {
	var dtos []proxy_models.UserDto
	for _, user := range users {
		dtos = append(dtos, mapToProxyDto(&user))
	}
	return dtos
}