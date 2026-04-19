package events

import (
	"time"

	"github.com/google/uuid"
)

type MessageSentEvent struct {
	Type      string
	ID        uuid.UUID
	ChatID    uuid.UUID
	SenderID  uuid.UUID
	Content   string
	Timestamp time.Time
}
