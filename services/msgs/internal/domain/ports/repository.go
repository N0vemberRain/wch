package ports

import (
	"context"
	"wch/services/msgs/internal/domain/model"

	"github.com/google/uuid"
)

type MessageRepository interface {
	Save(context.Context, *model.Message) error
	List(context.Context, uuid.UUID) ([]model.Message, error)
}
