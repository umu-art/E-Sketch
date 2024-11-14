package service

import (
	"est-proxy/src/config"
	"est-proxy/src/models"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository() *UserRepository {
	repositoryAddress := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", 
		config.POSTGRES_USERNAME, 
		config.POSTGRES_PASSWORD, 
		config.POSTGRES_HOST, 
		config.POSTGRES_PORT, 
		config.POSTGRES_DATABASE)

	config, err := pgxpool.ParseConfig(repositoryAddress)
	config.MaxConns = 20

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

func (r *UserRepository) Create(username string, email string, passwordHash string) error {
	_, err := r.db.Exec(context.Background(),  
		"INSERT INTO users (id, username, password_hash, email) VALUES ($1, $2, $3, $4)", 
		uuid.New(), username, passwordHash, email)
	if err != nil {
		return fmt.Errorf("Failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	var row pgx.Row

	row = r.db.QueryRow(context.Background(), 
		"SELECT id, username, password_hash, email, avatar FROM users WHERE email = $1", 
		email)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.Avatar)
	if err != nil {
		return nil, fmt.Errorf("Failed to find user with email %s: %w", email, err)
	}
	
	return &user, nil
}

func (r *UserRepository) GetIDByEmail(email string) (*uuid.UUID, error) {
	var id uuid.UUID
	var row pgx.Row

	row = r.db.QueryRow(context.Background(), 
		"SELECT id FROM users WHERE email = $1", 
		email)
	err := row.Scan(&id)
	if err == nil {
		return &id, nil
	}
	
	return nil, fmt.Errorf("Failed to find user with email %s: %w", email, err)
}

func (r *UserRepository) GetUserByID(id *uuid.UUID) (*models.User, error) {
	var user models.User
	var row pgx.Row

	row = r.db.QueryRow(context.Background(), 
		"SELECT id, username, password_hash, email, avatar FROM users WHERE id = $1", 
		id)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.Avatar)
	if err != nil {
		return nil, fmt.Errorf("Failed to find user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) UserExistsByUsernameOrEmail(username string, email string) (bool, error) {
	var count int

	err := r.db.QueryRow(context.Background(), 
		"SELECT COUNT(*) FROM users username = $1 OR email = $2", 
		username, email).Scan(&count)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	
	return count > 0, nil
}

func (r *UserRepository) SearchByUsernameIgnoreCase(ctx context.Context, username string) ([]models.User, error) {
	query := "SELECT id, username, password_hash, email, avatar FROM users WHERE username ILIKE '%' || $1 || '%'"
    rows, err := r.db.Query(ctx, query, username)
    if err != nil {
        return nil, fmt.Errorf("Failed to parse users, %w", err)
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var user models.User
        if err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.Avatar); err != nil {
            return nil, fmt.Errorf("Failed to parse users, %w", err)
        }
        users = append(users, user)
    }

    if rows.Err() != nil {
        return nil, fmt.Errorf("Failed to parse users, %w", rows.Err())
    }

    return users, nil
}