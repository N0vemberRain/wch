package ports

import (
	"context"

	"github.com/google/uuid"
)

type ChatProvider interface {
	GetUserIDsInChat(context.Context, uuid.UUID) (uuid.UUIDs, error)
}
