package ports

import (
	"context"
	"wch/services/auth/internal/domain/model"

	"github.com/google/uuid"
)

type CredentialsRepository interface {
	GetByEmail(context.Context, string) (*model.Credentials, error)
	GetByUserID(context.Context, uuid.UUID) (*model.Credentials, error)
	Save(context.Context, *model.Credentials) error
}
