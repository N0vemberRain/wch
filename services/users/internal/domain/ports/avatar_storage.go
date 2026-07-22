package repository

import (
	"context"

	"github.com/google/uuid"
)

type AvatarStorage interface {
	Save(ctx context.Context, userID uuid.UUID, data []byte) (string, error)
	Delete(ctx context.Context, key string) error
	Get(ctx context.Context, key string) ([]byte, error)
}
