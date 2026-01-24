package domain

import (
	"time"

	"github.com/google/uuid"
)

type Credentials struct {
	UserID       uuid.UUID
	Name         string
	PasswordHash string
	IsActive     bool
	CreatedAt    time.Time
}
