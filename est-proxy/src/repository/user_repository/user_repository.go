package user_repository

import (
	"context"
	"est-proxy/src/models"
	"est-proxy/src/repository/postgres"
	"est-proxy/src/repository/user_repository/impl"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, username string, email string, passwordHash string) *uuid.UUID
	GetUserByEmail(ctx context.Context, email string) *models.User
	GetIDByEmail(ctx context.Context, email string) *uuid.UUID
	GetUserByID(ctx context.Context, id *uuid.UUID) *models.User
	UserExistsByUsernameOrEmail(ctx context.Context, username string, email string) bool
	SearchByUsernameIgnoreCase(ctx context.Context, username string) *[]models.PublicUser
	GetUserListByIds(ctx context.Context, ids []uuid.UUID) *[]models.PublicUser
}

func NewUserRepository(postgresService postgres.PostgresService) UserRepository {
	return impl.NewUserRepositoryImpl(postgresService)
}
