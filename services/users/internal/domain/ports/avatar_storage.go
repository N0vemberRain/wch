package repository

import (
	"context"

	"github.com/google/uuid"
)

type AvatarStorage interface {
	Save(ctx context.Context, userID uuid.UUID, data []byte) (string, error)
	Delete(ctx context.Context, key string) error
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]byte, error)
	GetByKey(ctx context.Context, key string) ([]byte, error)
}
