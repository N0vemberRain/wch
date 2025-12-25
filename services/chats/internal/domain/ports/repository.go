package ports

import (
	"context"
	"wch/services/chats/internal/domain/model"

	"github.com/google/uuid"
)

type Repository interface {
	CreateChat(ctx context.Context, c *model.Chat) error
	GetChat(ctx context.Context, chatID uuid.UUID) (*model.Chat, error)
	UpdateChat(ctx context.Context, c *model.Chat) error
	DeleteChat(ctx context.Context, chatID uuid.UUID) error

	AddParticipant(ctx context.Context, p *model.ChatParticipant) error
	RemoveParticipant(ctx context.Context, chatID, userID uuid.UUID) error
	GetParticipantsIDs(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error)
}
