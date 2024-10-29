package models

import (
	"time"
	"github.com/google/uuid"
)

type UserToken struct {

	UserUUID uuid.UUID `json:"userUUID"`

	ExpirationTime time.Time `json:"expirationTime"`
}
