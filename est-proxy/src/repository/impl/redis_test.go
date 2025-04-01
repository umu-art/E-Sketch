package impl

import (
	"context"
	"est-proxy/src/models"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func newTestRedisClientImpl(t *testing.T, addr string) *RedisClientImpl {
	client := RedisClientImpl{
		redisURL:      addr,
		redisPassword: "",
		redisDB:       "0",
	}
	err := client.connect()
	assert.NoError(t, err)
	return &client
}

func TestRedisClientImpl_Normal(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	client := newTestRedisClientImpl(t, mr.Addr())
	ctx := context.Background()
	userKey := "user:123"
	user := &models.RegisteredUser{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	t.Run("AddUser", func(t *testing.T) {
		err = client.AddUser(ctx, userKey, user)
		assert.NoError(t, err)
	})

	t.Run("GetUser", func(t *testing.T) {
		val, err := mr.Get(userKey)
		assert.NoError(t, err)
		expectedValue := "testuser |#&^| test@example.com |#&^| hashedpassword"
		assert.Equal(t, expectedValue, val)

		retrievedUser, err := client.GetUser(ctx, userKey)
		assert.NoError(t, err)
		assert.Equal(t, user, retrievedUser)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		err = client.RemoveUser(ctx, userKey)
		assert.NoError(t, err)
		_, err = mr.Get(userKey)
		assert.Error(t, err)
	})
}

func TestRedisClientImpl_NonExistent(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	client := newTestRedisClientImpl(t, mr.Addr())
	ctx := context.Background()

	t.Run("GetUser", func(t *testing.T) {
		user, err := client.GetUser(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
	t.Run("RemoveUser", func(t *testing.T) {
		err = client.RemoveUser(ctx, "nonexistent")
		assert.NoError(t, err)
	})
}

func TestRedisClientImpl_Refresh(t *testing.T) {
	originalRefreshPeriod := refreshPeriod
	defer func() {
		refreshPeriod = originalRefreshPeriod
	}()
	refreshPeriod = 5 * time.Millisecond

	mr, err := miniredis.Run()
	assert.NoError(t, err)

	client := newTestRedisClientImpl(t, mr.Addr())

	go client.Refresh()

	ctx := context.Background()
	userKey := "user1"
	user := &models.RegisteredUser{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	err = client.AddUser(ctx, userKey, user)
	assert.NoError(t, err)

	mr.Close()

	time.Sleep(refreshPeriod)

	mr, err = miniredis.Run()
	assert.NoError(t, err)
	client.redisURL = mr.Addr()

	time.Sleep(time.Second)

	err = client.AddUser(ctx, "user2", user)
	assert.NoError(t, err)
}
