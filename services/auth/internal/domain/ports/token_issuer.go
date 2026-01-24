package ports

import (
	"wch/services/auth/internal/domain"

	"github.com/google/uuid"
)

type TokenIssuer interface {
	Issue(userID uuid.UUID) (domain.Token, error)
}
