package ports

import (
	"context"
	"wch/services/chats/internal/domain/model"

	"github.com/google/uuid"
)

type UserProvider interface {
	GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]model.User, error)
	GetAvatarsForChats(ctx context.Context, ids []uuid.UUID) ([]model.Avatar, error)
	UpdateAvatarForChat(ctx context.Context, av *model.Avatar) (string, error)
}
