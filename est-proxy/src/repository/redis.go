package repository

import (
	"context"
	"est-proxy/src/models"
)

type RedisClient interface {
	AddUser(ctx context.Context, userKey string, user *models.RegisteredUser) error
	GetUser(ctx context.Context, userKey string) (*models.RegisteredUser, error)
	RemoveUser(ctx context.Context, userKey string) error
	Refresh()
	Close()
}
