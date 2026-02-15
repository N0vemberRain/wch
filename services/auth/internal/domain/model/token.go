package model

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	Value     string
	ExpiresAt time.Time
	UserID    uuid.UUID
}
