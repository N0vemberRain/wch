package model

import (
	"time"
	"wch/services/chats"

	"github.com/google/uuid"
)

type ChatType int32

const (
	ChatTypeUnknown ChatType = 0
	ChatTypeDirect  ChatType = 1
	ChatTypeGroup   ChatType = 2
)

type Chat struct {
	ID        uuid.UUID
	Type      ChatType
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewChat(
	id uuid.UUID,
	chatType ChatType,
	name string,
	createdAt time.Time,
	updatedAt time.Time,
) (*Chat, error) {
	if id == uuid.Nil {
		return nil, chats.ErrChatIDNil
	}

	return &Chat{
		ID:        id,
		Type:      chatType,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
