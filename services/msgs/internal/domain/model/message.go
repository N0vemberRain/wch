package model

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Id        uuid.UUID
	ChatId    uuid.UUID
	SenderId  uuid.UUID
	Content   string
	CreatedAt time.Time
}
