package model

import (
	"time"
	"wch/services/chats"

	"github.com/google/uuid"
)

type ChatParticipantRole int32

const (
	ChatParticipantAdmin   = 1
	ChatParticipantMember  = 2
	ChatParticipantUnknown = 0
)

type ChatParticipant struct {
	UserID   uuid.UUID
	ChatID   uuid.UUID
	Role     ChatParticipantRole
	JoinedAt time.Time
}

func NewChatParticipant(
	userID uuid.UUID,
	chatID uuid.UUID,
	role ChatParticipantRole,
	joinedAt time.Time,
) (*ChatParticipant, error) {
	if userID == uuid.Nil {
		return nil, chats.ErrUserIDNil
	}

	if chatID == uuid.Nil {
		return nil, chats.ErrChatIDNil
	}

	return &ChatParticipant{
		UserID:   userID,
		ChatID:   chatID,
		Role:     role,
		JoinedAt: joinedAt,
	}, nil
}
