package ports

import (
	"context"
	"wch/services/msgs/internal/domain/model"

	"github.com/google/uuid"
)

type ChatProvider interface {
	IsUserParticipant(context.Context, uuid.UUID, uuid.UUID) (bool, model.ChatParticipantRole, error)
	ChatExists(context.Context, uuid.UUID) (bool, error)
}
