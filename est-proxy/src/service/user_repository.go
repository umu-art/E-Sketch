package service

import (
	"context"
	"est-proxy/src/models"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository() *UserRepository {
	config, err := pgxpool.ParseConfig(os.Getenv("USER_REPOSIRORY_ADDRESS"))
    if err != nil {
        log.Fatalf("Unable to parse config: %v", err)
    }

    db, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }

    return &UserRepository{db: db}
}

func (r *UserRepository) Release() {
	if r.db != nil {
		r.db.Close()
		r.db = nil
	}
}

func (r *UserRepository) Create(username string, email string, password_hash string) error {
	if r.db == nil {
		return fmt.Errorf("Unverified query to closed db")
	}
	_, err := r.db.Exec(context.Background(),  
		"INSERT INTO users (id, username, password_hash, email) VALUES ($1, $2, $3, $4)", 
		uuid.New(), username, password_hash, email)
	if (err != nil) {
		return fmt.Errorf("Failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetIDByUsername(username string) (*uuid.UUID, error) {
	var id uuid.UUID
	var row pgx.Row
	var err error = nil

	if r.db == nil {
		return nil, fmt.Errorf("Unverified query to closed db")
	}

	row = r.db.QueryRow(context.Background(), "SELECT id FROM users WHERE username = $1", username)
	err = row.Scan(&id)
	if err == nil {
		return &id, nil
	}
	
	return nil, fmt.Errorf("Failed to find user with username %s: %w", username, err)
}

func (r *UserRepository) GetIDByEmail(email string) (*uuid.UUID, error) {
	var id uuid.UUID
	var row pgx.Row
	var err error = nil

	if r.db == nil {
		return nil, fmt.Errorf("Unverified query to closed db")
	}

	row = r.db.QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1", email)
	err = row.Scan(&id)
	if err == nil {
		return &id, nil
	}
	
	return nil, fmt.Errorf("Failed to find user with email %s: %w", email, err)
}

func (r *UserRepository) GetUserByID(id *uuid.UUID) (*models.User, error) {
	var user models.User
	var row pgx.Row

	if r.db == nil {
		return nil, fmt.Errorf("Unverified query to closed db")
	}

	row = r.db.QueryRow(context.Background(), "SELECT id, username, password_hash, email, avatar FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.Avatar)
	if err != nil {
		return nil, fmt.Errorf("Failed to find user: %w", err)
	}

	return &user, nil
}