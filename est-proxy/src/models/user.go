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

	IsBanned bool
}

type PublicUser struct {
	ID uuid.UUID

	Username string

	Avatar string
}

type RegisteredUser struct {
	Username string

	PasswordHash string

	Email string
}

func (u User) Public() *PublicUser {
	return &PublicUser{
		ID:       u.ID,
		Username: u.Username,
		Avatar:   u.Avatar,
	}
}
