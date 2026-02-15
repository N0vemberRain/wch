package ports

import (
	"wch/services/auth/internal/domain/model"

	"github.com/google/uuid"
)

type TokenIssuer interface {
	Issue(userID uuid.UUID) (*model.Token, error)
	Validate(token string) (*model.Token, error)
}
