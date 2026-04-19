package domain

import "github.com/google/uuid"

type Notification struct {
	UserID  uuid.UUID
	ChatID  uuid.UUID
	Content string
}
