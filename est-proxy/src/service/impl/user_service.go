package impl

import (
	"context"
	"crypto/rand"
	"est-proxy/src/api"
	"est-proxy/src/config"
	"est-proxy/src/errors"
	"est-proxy/src/mapper"
	"est-proxy/src/repository"
	"est-proxy/src/utils"
	proxymodels "est_proxy_go/models"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type UserServiceImpl struct {
	userRepository repository.UserRepository
	mailApi        *api.MailApi
	redisClient    repository.RedisClient
}

func NewUserServiceImpl(userRepository repository.UserRepository, mailApi *api.MailApi, redisClient repository.RedisClient) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
		mailApi:        mailApi,
		redisClient:    redisClient,
	}
}

func (u UserServiceImpl) GetUserById(ctx context.Context, userId *uuid.UUID) (*proxymodels.UserDto, *errors.StatusError) {
	user := u.userRepository.GetUserByID(ctx, userId)
	if user == nil {
		return nil, errors.NewStatusError(http.StatusBadRequest, "Пользователь не найден")
	}

	return mapper.MapUserToProxy(*user.Public()), nil
}

func (u UserServiceImpl) Login(ctx context.Context, authDto *proxymodels.AuthDto) (*string, *errors.StatusError) {
	user := u.userRepository.GetUserByEmail(ctx, authDto.Email)
	if user == nil {
		return nil, errors.NewStatusError(http.StatusBadRequest, "Отсутствует или некорректный адрес почты или пароль")
	}

	if user.IsBanned {
		return nil, errors.NewStatusError(http.StatusForbidden, "Аккаунт заблокирован")
	}

	if authDto.PasswordHash != user.PasswordHash {
		return nil, errors.NewStatusError(http.StatusBadRequest, "Отсутствует или некорректный адрес почты или пароль")
	}

	u.userRepository.UpdateLoggedInUser(ctx, &user.ID)

	token := utils.GenerateUserJWTString(&user.ID)
	if token == nil {
		return nil, errors.NewStatusError(http.StatusUnauthorized, "Не получилось войти в аккаунт")
	}

	return token, nil
}

func (u UserServiceImpl) Register(ctx context.Context, registerDto *proxymodels.RegisterDto) *errors.StatusError {
	exists := u.userRepository.UserExistsByUsernameOrEmail(ctx, registerDto.Username, registerDto.Email)
	if exists {
		return errors.NewStatusError(http.StatusConflict, "Занято имя пользователя или адрес электронной почты")
	}

	user := mapper.MapProxyRegisterDto(registerDto)

	token, err := u.generateToken(10)
	if err != nil {
		return errors.NewStatusError(http.StatusInternalServerError, "Не получилось отправить письмо для подтверждения почты")
	}
	confirmationLink := u.generateConfirmationLink(token)

	err = u.redisClient.AddUser(ctx, token, user)
	if err != nil {
		return errors.NewStatusError(http.StatusInternalServerError, "Не получилось отправить письмо для подтверждения почты")
	}

	err = u.mailApi.SendConfirmationEmail(user.Email, confirmationLink)
	if err != nil {
		return errors.NewStatusError(http.StatusInternalServerError, "Не получилось отправить письмо для подтверждения почты")
	}

	return nil
}

func (u UserServiceImpl) Confirm(ctx context.Context, userToken string) (*string, *errors.StatusError) {
	user, err := u.redisClient.GetUser(ctx, userToken)
	if err != nil {
		return nil, errors.NewStatusError(http.StatusUnauthorized, "Проверьте правильность кода")
	}

	userId := u.userRepository.Create(ctx, user.Username, user.Email, user.PasswordHash)
	if userId == nil {
		return nil, errors.NewStatusError(http.StatusUnauthorized, "Не получилось подтвердить аккаунт")
	}

	token := utils.GenerateUserJWTString(userId)
	if token == nil {
		return nil, errors.NewStatusError(http.StatusUnauthorized, "Не получилось войти в аккаунт")
	}

	err = u.redisClient.RemoveUser(ctx, userToken)
	if err != nil {
		log.Println(err.Error())
	}

	return token, nil
}

func (u UserServiceImpl) Search(ctx context.Context, query string) (*[]proxymodels.UserDto, *errors.StatusError) {
	users := u.userRepository.SearchByUsernameIgnoreCase(ctx, query)
	if users == nil {
		return nil, errors.NewStatusError(http.StatusInternalServerError, "Ошибка поиска")
	}

	return mapper.MapUserListToProxy(*users), nil
}

func (u UserServiceImpl) generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		char := bytes[i] % 62
		if char < 10 {
			result[i] = char + '0'
		} else if char < 36 {
			result[i] = char - 10 + 'A'
		} else {
			result[i] = char - 36 + 'a'
		}
	}

	log.Printf("New user token: %s", string(result))

	return string(result), nil
}

func (u UserServiceImpl) generateConfirmationLink(token string) string {
	return fmt.Sprintf("%s?token=%s", config.CONFIRM_URL, token)
}
