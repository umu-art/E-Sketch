package models

import (
	"time"

	"github.com/google/uuid"
)

type ParsedJWT struct {
	UserID uuid.UUID

	ExpirationTime time.Time
}
