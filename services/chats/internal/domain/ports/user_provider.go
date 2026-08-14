package ports

import (
	"context"
	"wch/services/chats/internal/domain/model"

	"github.com/google/uuid"
)

type UserProvider interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]model.User, error)
	GetAvatarForDirectChatByKey(ctx context.Context, key string) (*model.Avatar, error)
	GetAvatarsForChats(ctx context.Context, ids []uuid.UUID) ([]model.Avatar, error)
	UpdateAvatarForChat(ctx context.Context, av *model.Avatar) (string, error)
}
