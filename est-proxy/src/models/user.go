package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID

	Username string

	PasswordHash string

	Email string

	Avatar string
}
