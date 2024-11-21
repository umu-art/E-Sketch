package service

import (
	"context"
	"est-proxy/src/config"
	"est-proxy/src/models"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/url"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository() *UserRepository {
	repositoryAddress := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.POSTGRES_USERNAME,
		url.QueryEscape(config.POSTGRES_PASSWORD),
		config.POSTGRES_HOST,
		config.POSTGRES_PORT,
		config.POSTGRES_DATABASE)

	pgxConfig, err := pgxpool.ParseConfig(repositoryAddress)
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}

	pgxConfig.MaxConns = 20

	db, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	return &UserRepository{db: db}
}

func (r UserRepository) Release() {
	if r.db != nil {
		r.db.Close()
		r.db = nil
	}
}

func (r UserRepository) Create(username string, email string, passwordHash string) *uuid.UUID {
	_, err := r.db.Exec(context.Background(),
		"INSERT INTO users (id, username, password_hash, email) VALUES ($1, $2, $3, $4)",
		uuid.New(), username, passwordHash, email)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil
	}
	return r.GetIDByEmail(email)
}

func (r UserRepository) GetUserByEmail(email string) *models.User {
	var user models.User
	var row pgx.Row

	row = r.db.QueryRow(context.Background(),
		"SELECT id, username, password_hash, email, avatar FROM users WHERE email = $1",
		email)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.Avatar)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return nil
	}

	return &user
}

func (r UserRepository) GetIDByEmail(email string) *uuid.UUID {
	var id uuid.UUID
	var row pgx.Row

	row = r.db.QueryRow(context.Background(),
		"SELECT id FROM users WHERE email = $1",
		email)
	err := row.Scan(&id)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return nil
	}

	return &id
}

func (r UserRepository) GetUserByID(id *uuid.UUID) *models.User {
	var user models.User
	var row pgx.Row

	row = r.db.QueryRow(context.Background(),
		"SELECT id, username, email, avatar FROM users WHERE id = $1",
		id)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Avatar)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return nil
	}

	return &user
}

func (r UserRepository) UserExistsByUsernameOrEmail(username string, email string) bool {
	var count int

	err := r.db.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2",
		username, email).Scan(&count)

	if err != nil {
		log.Printf("Failed to count users: %v", err)
		return true
	}

	return count > 0
}

func (r UserRepository) SearchByUsernameIgnoreCase(ctx context.Context, username string) *[]models.PublicUser {
	query := "SELECT id, username, avatar FROM users WHERE username ILIKE '%' || $1 || '%'"
	rows, err := r.db.Query(ctx, query, username)
	if err != nil {
		log.Printf("Failed to search users: %v", err)
		return nil
	}
	defer rows.Close()

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
