package ports

import (
	"context"
	"wch/services/auth/internal/domain"

	"github.com/google/uuid"
)

type CredentialsRepository interface {
	GetByEmail(context.Context, string) (*domain.Credentials, error)
	GetBuUserID(context.Context, uuid.UUID) (*domain.Credentials, error)
	Save(context.Context, *domain.Credentials) error
}
