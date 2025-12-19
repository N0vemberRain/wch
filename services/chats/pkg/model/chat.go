package model

import (
	"time"

	"github.com/google/uuid"
)

type ChatType int32

const (
	ChatTypeUnknown ChatType = 0
	ChatTypeDirect  ChatType = 1
	ChatTypeGroup   ChatType = 2
)

type ChatOption func(*Chat)

func WithID(id uuid.UUID) ChatOption {
	return func(c *Chat) {
		c.ID = id
	}
}

func WithCreatedTime(t time.Time) ChatOption {
	return func(c *Chat) {
		c.CreatedAt = t
	}
}

func WithUpdatedTime(t time.Time) ChatOption {
	return func(c *Chat) {
		c.UpdatedAt = t
	}
}

type Chat struct {
	ID        uuid.UUID
	Type      ChatType
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewChat(name string, t ChatType, opts ...ChatOption) *Chat {
	chat := &Chat{
		Name: name,
		Type: t,
	}

	for _, opt := range opts {
		opt(chat)
	}

	return chat
}

// func NewChat(
// 	id uuid.UUID,
// 	chatType ChatType,
// 	name string,
// 	createdAt time.Time,
// 	updatedAt time.Time,
// ) (*Chat, error) {
// 	if id == uuid.Nil {
// 		return nil, chats.ErrChatIDNil
// 	}

// 	return &Chat{
// 		ID:        id,
// 		Type:      chatType,
// 		Name:      name,
// 		CreatedAt: createdAt,
// 		UpdatedAt: updatedAt,
// 	}, nil
// }
