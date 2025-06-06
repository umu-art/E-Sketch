package impl

import (
	"context"
	"est-proxy/src/models"
	"est-proxy/src/repository"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

type UserRepositoryImpl struct {
	postgresService repository.PostgresService
}

func NewUserRepositoryImpl(postgresService repository.PostgresService) *UserRepositoryImpl {
	return &UserRepositoryImpl{postgresService}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, username string, email string, passwordHash string) *uuid.UUID {
	_, err := r.postgresService.Exec(ctx,
		"INSERT INTO users (id, username, password_hash, email, avatar, last_login) VALUES ($1, $2, $3, $4, '', $5)",
		uuid.New(), username, passwordHash, email, time.Now())
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil
	}
	return r.GetIDByEmail(ctx, email)
}

func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) *models.User {
	var user models.User
	var row pgx.Row

	row = r.postgresService.QueryRow(ctx,
		"SELECT id, username, password_hash, email, avatar, is_banned FROM users WHERE email = $1",
		email)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.Avatar, &user.IsBanned)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return nil
	}

	return &user
}

func (r *UserRepositoryImpl) GetIDByEmail(ctx context.Context, email string) *uuid.UUID {
	var id uuid.UUID
	var row pgx.Row

	row = r.postgresService.QueryRow(ctx,
		"SELECT id FROM users WHERE email = $1",
		email)
	err := row.Scan(&id)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return nil
	}

	return &id
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, id *uuid.UUID) *models.User {
	var user models.User
	var row pgx.Row

	row = r.postgresService.QueryRow(ctx,
		"SELECT id, username, email, avatar, is_banned FROM users WHERE id = $1",
		id)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Avatar, &user.IsBanned)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return nil
	}

	return &user
}

func (r *UserRepositoryImpl) UserExistsByUsernameOrEmail(ctx context.Context, username string, email string) bool {
	var count int

	err := r.postgresService.QueryRow(ctx,
		"SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2",
		username, email).Scan(&count)

	if err != nil {
		log.Printf("Failed to count users: %v", err)
		return true
	}

	return count > 0
}

func (r *UserRepositoryImpl) SearchByUsernameIgnoreCase(ctx context.Context, username string) *[]models.PublicUser {
	rows, err := r.postgresService.Query(ctx,
		"SELECT id, username, avatar FROM users WHERE username ILIKE '%' || $1 || '%'",
		username)
	defer rows.Close()
	if err != nil {
		log.Printf("Failed to search users: %v", err)
		return nil
	}

	var users []models.PublicUser
	for rows.Next() {
		var user models.PublicUser
		if err := rows.Scan(&user.ID, &user.Username, &user.Avatar); err != nil {
			fmt.Printf("Failed to parse user, %v \n", err)
			continue
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		log.Printf("Failed to search users: %v", rows.Err())
		return nil
	}

	return &users
}

func (r *UserRepositoryImpl) GetUserListByIds(ctx context.Context, ids []uuid.UUID) *[]models.PublicUser {
	rows, err := r.postgresService.Query(ctx,
		"SELECT id, username, avatar FROM users WHERE id = ANY($1)",
		ids)
	defer rows.Close()
	if err != nil {
		log.Printf("Failed to search users: %v", err)
		return nil
	}

	var users []models.PublicUser
	for rows.Next() {
		var user models.PublicUser
		if err := rows.Scan(&user.ID, &user.Username, &user.Avatar); err != nil {
			fmt.Printf("Failed to parse user, %v \n", err)
			continue
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		log.Printf("Failed to search users: %v", rows.Err())
		return nil
	}

	return &users
}

func (r *UserRepositoryImpl) UpdateLoggedInUser(ctx context.Context, userId *uuid.UUID) {
	_, err := r.postgresService.Exec(ctx, "UPDATE users SET last_login = $1 WHERE id = $2", time.Now(), userId)
	if err != nil {
		log.Printf("Failed to update logged in user: %v", err)
	}
}
