package service

import (
	"context"
	"est-proxy/src/models"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// users: [id, username, password_hash, email, avatar]

type UserRepository struct {
	db *pgxpool.Pool
	isOpened bool
}

func NewUserRepository(dbconfig string) (*UserRepository) {
	config, err := pgxpool.ParseConfig(dbconfig)
    if err != nil {
        log.Fatalf("Unable to parse config: %v", err)
    }

    db, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }

    return &UserRepository{db: db, isOpened: true}
}

func (r *UserRepository) Release() {
	if r.isOpened {
		r.db.Close()
		r.isOpened = false
	}
}

func (r *UserRepository) Create(username string, email string, password_hash string) error {
	if !r.isOpened {
		return fmt.Errorf("unverified query to closed db")
	}
	_, err := r.db.Exec(context.Background(),  
		"INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3)", 
		username, password_hash, email)
	if (err != nil) {
		return fmt.Errorf("(postgress) failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetIDByUserInfo(username *string, email *string) (*uuid.UUID, error) {
	var id uuid.UUID
	var row pgx.Row
	var err error = nil

	if !r.isOpened {
		return nil, fmt.Errorf("unverified query to closed db")
	}

	if username == nil && email == nil {
		return nil, fmt.Errorf("no user info provided")
	}

	if username != nil {
		row = r.db.QueryRow(context.Background(), "SELECT id FROM users WHERE username = $1", username)
		err = row.Scan(&id)
		if err == nil {
			return &id, nil
		}
	}
	if email != nil {
		row = r.db.QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1", email)
		err = row.Scan(&id)
		if err == nil {
			return &id, nil
		}
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	
}

func (r *UserRepository) GetUserByID(id *uuid.UUID) (*models.User, error) {
	var user models.User
	var row pgx.Row

	if !r.isOpened {
		return nil, fmt.Errorf("unverified query to closed db")
	}

	row = r.db.QueryRow(context.Background(), "SELECT id, username, password_hash, email, avatar FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.Avatar)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

