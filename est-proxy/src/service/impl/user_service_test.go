package impl

import (
	"context"
	"errors"
	"est-proxy/src/models"
	"fmt"
	"net/http"
	"testing"

	"est-proxy/src/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	proxymodels "est_proxy_go/models"
)

func TestUserServiceImpl_GetUserById(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockMail := new(MockMailApi)
	mockRedis := new(MockRedisClient)
	service := NewUserServiceImpl(mockRepo, mockMail, mockRedis)

	t.Run("user exists", func(t *testing.T) {
		user := newDummyUser()
		userID := user.ID
		mockRepo.
			On("GetUserByID", ctx, &userID).
			Return(user).
			Once()

		result, errStatus := service.GetUserById(ctx, &userID)
		assert.NotNil(t, result)
		assert.Nil(t, errStatus)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		uid := uuid.New()
		mockRepo.
			On("GetUserByID", ctx, &uid).
			Return(nil).
			Once()

		result, errStatus := service.GetUserById(ctx, &uid)
		assert.Nil(t, result)
		assert.NotNil(t, errStatus)
		assert.Equal(t, http.StatusBadRequest, errStatus.HttpStatus)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserServiceImpl_Login(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockMail := new(MockMailApi)
	mockRedis := new(MockRedisClient)
	service := NewUserServiceImpl(mockRepo, mockMail, mockRedis)
	dummyUser := newDummyUser()

	t.Run("user not found", func(t *testing.T) {
		email := "nonexistent@example.com"
		authDto := &proxymodels.AuthDto{
			Email:        email,
			PasswordHash: "any",
		}
		mockRepo.
			On("GetUserByEmail", ctx, email).
			Return(nil).
			Once()

		token, errStatus := service.Login(ctx, authDto)
		assert.Nil(t, token)
		assert.NotNil(t, errStatus)
		assert.Equal(t, http.StatusBadRequest, errStatus.HttpStatus)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user banned", func(t *testing.T) {
		bannedUser := newDummyUser()
		bannedUser.IsBanned = true
		authDto := &proxymodels.AuthDto{
			Email:        bannedUser.Email,
			PasswordHash: bannedUser.PasswordHash,
		}
		mockRepo.
			On("GetUserByEmail", ctx, bannedUser.Email).
			Return(bannedUser).
			Once()

		token, errStatus := service.Login(ctx, authDto)
		assert.Nil(t, token)
		assert.NotNil(t, errStatus)
		assert.Equal(t, http.StatusForbidden, errStatus.HttpStatus)
		mockRepo.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		authDto := &proxymodels.AuthDto{
			Email:        dummyUser.Email,
			PasswordHash: "wronghash",
		}
		mockRepo.
			On("GetUserByEmail", ctx, dummyUser.Email).
			Return(dummyUser).
			Once()

		token, errStatus := service.Login(ctx, authDto)
		assert.Nil(t, token)
		assert.NotNil(t, errStatus)
		assert.Equal(t, http.StatusBadRequest, errStatus.HttpStatus)
		mockRepo.AssertExpectations(t)
	})

	t.Run("successful login", func(t *testing.T) {
		authDto := &proxymodels.AuthDto{
			Email:        dummyUser.Email,
			PasswordHash: dummyUser.PasswordHash,
		}
		mockRepo.
			On("GetUserByEmail", ctx, dummyUser.Email).
			Return(dummyUser).
			Once()
		mockRepo.
			On("UpdateLoggedInUser", ctx, &dummyUser.ID).
			Return().
			Once()

		generateJWTString = func(userID *uuid.UUID) *string {
			token := "dummy-token"
			return &token
		}

		token, errStatus := service.Login(ctx, authDto)
		assert.NotNil(t, token)
		assert.Equal(t, "dummy-token", *token)
		assert.Nil(t, errStatus)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserServiceImpl_Register(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockMail := new(MockMailApi)
	mockRedis := new(MockRedisClient)
	service := NewUserServiceImpl(mockRepo, mockMail, mockRedis)

	dummyRegister := &proxymodels.RegisterDto{
		Username:     "newuser",
		Email:        "new@example.com",
		PasswordHash: "newhash",
	}

	t.Run("user already exists", func(t *testing.T) {
		mockRepo.
			On("UserExistsByUsernameOrEmail", ctx, dummyRegister.Username, dummyRegister.Email).
			Return(true).
			Once()

		errStatus := service.Register(ctx, dummyRegister)
		assert.NotNil(t, errStatus)
		assert.Equal(t, http.StatusConflict, errStatus.HttpStatus)
		mockRepo.AssertExpectations(t)
	})

	t.Run("redis add failure", func(t *testing.T) {
		mockRepo.
			On("UserExistsByUsernameOrEmail", ctx, dummyRegister.Username, dummyRegister.Email).
			Return(false).
			Once()

		mockRedis.
			On("AddUser", ctx, mock.Anything, mock.Anything).
			Return(errors.New("redis error")).
			Once()

		errStatus := service.Register(ctx, dummyRegister)
		assert.NotNil(t, errStatus)
		assert.Equal(t, http.StatusInternalServerError, errStatus.HttpStatus)
		mockRedis.AssertExpectations(t)
	})

	t.Run("mail sending failure", func(t *testing.T) {
		mockRepo.
			On("UserExistsByUsernameOrEmail", ctx, dummyRegister.Username, dummyRegister.Email).
			Return(false).
			Once()
		mockRedis.
			On("AddUser", ctx, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		mockMail.
			On("SendConfirmationEmail", dummyRegister.Email, mock.Anything, mock.Anything).
			Return(errors.New("mail error")).
			Once()

		errStatus := service.Register(ctx, dummyRegister)
		assert.NotNil(t, errStatus)
		assert.Equal(t, http.StatusInternalServerError, errStatus.HttpStatus)
		mockRepo.AssertExpectations(t)
		mockMail.AssertExpectations(t)
		mockRedis.AssertExpectations(t)
	})

	t.Run("successful registration", func(t *testing.T) {
		mockRepo.
			On("UserExistsByUsernameOrEmail", ctx, dummyRegister.Username, dummyRegister.Email).
			Return(false).
			Once()
		mockRedis.
			On("AddUser", ctx, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		mockMail.
			On("SendConfirmationEmail", dummyRegister.Email, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		errStatus := service.Register(ctx, dummyRegister)
		assert.Nil(t, errStatus)
		mockRepo.AssertExpectations(t)
		mockRedis.AssertExpectations(t)
		mockMail.AssertExpectations(t)
	})
}

func TestUserServiceImpl_Confirm(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockMail := new(MockMailApi)
	mockRedis := new(MockRedisClient)
	service := NewUserServiceImpl(mockRepo, mockMail, mockRedis)

	dummyToken := "dummytoken"
	dummyRegister := &models.RegisteredUser{
		Username:     "newuser",
		Email:        "new@example.com",
		PasswordHash: "newhash",
	}
	dummyId := uuid.New()

	defer func() {
		generateJWTString = originalGenerateJWTString
	}()

	t.Run("invalid token", func(t *testing.T) {
		mockRedis.
			On("GetUser", ctx, dummyToken).
			Return(nil, errors.New("not found")).
			Once()

		token, errStatus := service.Confirm(ctx, dummyToken)
		assert.Nil(t, token)
		assert.NotNil(t, errStatus)
		assert.Equal(t, http.StatusUnauthorized, errStatus.HttpStatus)
		mockRedis.AssertExpectations(t)
	})

	t.Run("user creation failure", func(t *testing.T) {
		mockRedis.
			On("GetUser", ctx, dummyToken).
			Return(dummyRegister, nil).
			Once()
		mockRepo.
			On("Create", ctx, dummyRegister.Username, dummyRegister.Email, dummyRegister.PasswordHash).
			Return(nil).
			Once()

		token, errStatus := service.Confirm(ctx, dummyToken)
		assert.Nil(t, token)
		assert.NotNil(t, errStatus)
		assert.Equal(t, http.StatusUnauthorized, errStatus.HttpStatus)
		mockRedis.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("JWT generation failure", func(t *testing.T) {
		mockRedis.
			On("GetUser", ctx, dummyToken).
			Return(dummyRegister, nil).
			Once()
		mockRepo.
			On("Create", ctx, dummyRegister.Username, dummyRegister.Email, dummyRegister.PasswordHash).
			Return(&dummyId).
			Once()
		generateJWTString = func(userID *uuid.UUID) *string {
			return nil
		}

		token, errStatus := service.Confirm(ctx, dummyToken)
		assert.Nil(t, token)
		assert.NotNil(t, errStatus)
		assert.Equal(t, http.StatusUnauthorized, errStatus.HttpStatus)
		mockRedis.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("successful confirmation", func(t *testing.T) {
		mockRedis.
			On("GetUser", ctx, dummyToken).
			Return(dummyRegister, nil).
			Once()
		mockRepo.
			On("Create", ctx, dummyRegister.Username, dummyRegister.Email, dummyRegister.PasswordHash).
			Return(&dummyId).
			Once()

		generateJWTString = func(userID *uuid.UUID) *string {
			token := "confirmed-token"
			return &token
		}
		mockRedis.
			On("RemoveUser", ctx, dummyToken).
			Return(nil).
			Once()

		token, errStatus := service.Confirm(ctx, dummyToken)
		assert.NotNil(t, token)
		assert.Equal(t, "confirmed-token", *token)
		assert.Nil(t, errStatus)
		mockRedis.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserServiceImpl_Search(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockMail := new(MockMailApi)
	mockRedis := new(MockRedisClient)
	service := NewUserServiceImpl(mockRepo, mockMail, mockRedis)

	t.Run("search failure", func(t *testing.T) {
		query := "nonexistent"
		mockRepo.
			On("SearchByUsernameIgnoreCase", ctx, query).
			Return(nil).
			Once()

		result, errStatus := service.Search(ctx, query)
		assert.Nil(t, result)
		assert.NotNil(t, errStatus)
		assert.Equal(t, http.StatusInternalServerError, errStatus.HttpStatus)
		mockRepo.AssertExpectations(t)
	})

	t.Run("successful search", func(t *testing.T) {
		query := "test"
		dummyUser := newDummyUser()
		userList := []models.PublicUser{*dummyUser.Public()}
		mockRepo.
			On("SearchByUsernameIgnoreCase", ctx, query).
			Return(&userList).
			Once()

		result, errStatus := service.Search(ctx, query)
		assert.NotNil(t, result)
		assert.Nil(t, errStatus)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserServiceImpl_generateConfirmationLink(t *testing.T) {
	service := NewUserServiceImpl(nil, nil, nil)
	token := "testtoken"
	expected := fmt.Sprintf("%s?token=%s", config.CONFIRM_URL, token)
	result := service.generateConfirmationLink(token)
	assert.Equal(t, expected, result)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetIDByEmail(ctx context.Context, email string) *uuid.UUID {
	return nil
}

func (m *MockUserRepository) GetUserListByIds(ctx context.Context, ids []uuid.UUID) *[]models.PublicUser {
	return nil
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id *uuid.UUID) *models.User {
	args := m.Called(ctx, id)
	if user, ok := args.Get(0).(*models.User); ok {
		return user
	}
	return nil
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) *models.User {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(*models.User); ok {
		return user
	}
	return nil
}

func (m *MockUserRepository) UpdateLoggedInUser(ctx context.Context, id *uuid.UUID) {
	m.Called(ctx, id)
}

func (m *MockUserRepository) UserExistsByUsernameOrEmail(ctx context.Context, username, email string) bool {
	args := m.Called(ctx, username, email)
	return args.Bool(0)
}

func (m *MockUserRepository) Create(ctx context.Context, username, email, passwordHash string) *uuid.UUID {
	args := m.Called(ctx, username, email, passwordHash)
	if id, ok := args.Get(0).(*uuid.UUID); ok {
		return id
	}
	return nil
}

func (m *MockUserRepository) SearchByUsernameIgnoreCase(ctx context.Context, query string) *[]models.PublicUser {
	args := m.Called(ctx, query)
	if users, ok := args.Get(0).(*[]models.PublicUser); ok {
		return users
	}
	return nil
}

type MockMailApi struct {
	mock.Mock
}

func (m *MockMailApi) SendConfirmationEmail(email, confirmationLink, token string) error {
	args := m.Called(email, confirmationLink, token)
	return args.Error(0)
}

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) AddUser(ctx context.Context, userKey string, user *models.RegisteredUser) error {
	args := m.Called(ctx, userKey, *user)
	return args.Error(0)
}

func (m *MockRedisClient) GetUser(ctx context.Context, userKey string) (*models.RegisteredUser, error) {
	args := m.Called(ctx, userKey)
	if user, ok := args.Get(0).(*models.RegisteredUser); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRedisClient) RemoveUser(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRedisClient) Refresh() {}
func (m *MockRedisClient) Close()   {}

func newDummyUser() *models.User {
	return &models.User{
		ID:           uuid.New(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash123",
		IsBanned:     false,
	}
}

var originalGenerateJWTString = generateJWTString
