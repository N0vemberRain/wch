package controller

import (
	"context"
	"wch/services/msgs/internal/domain"
	"wch/services/msgs/internal/domain/model"
	"wch/services/msgs/internal/domain/ports"

	"github.com/google/uuid"
)

type Controller struct {
	repo ports.MessageRepository
}

func NewMessageController(repo ports.MessageRepository) *Controller {
	return &Controller{
		repo: repo,
	}
}

func (c *Controller) Save(ctx context.Context, msg *model.Message) error {
	msg.Id = uuid.New()
	return c.repo.Save(ctx, msg)
}

func (c *Controller) List(ctx context.Context, chat_id uuid.UUID) ([]model.Message, error) {
	if chat_id == uuid.Nil {
		return nil, domain.ErrChatIsEmpty
	}

	return c.repo.List(ctx, chat_id)
}
