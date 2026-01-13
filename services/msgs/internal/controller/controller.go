package controller

import (
	"context"
	"time"
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
	msg.CreatedAt = time.Now()
	return c.repo.Save(ctx, msg)
}

func (c *Controller) List(ctx context.Context, chat_id uuid.UUID, limit int, cursor *model.Cursor) ([]model.Message, string, error) {
	if chat_id == uuid.Nil {
		return nil, "", domain.ErrChatIsEmpty
	}

	msgs, err := c.repo.List(ctx, chat_id, limit, cursor)
	var nextCursor string
	if err == nil && len(msgs) >= limit {
		nextCursor = msgs[len(msgs)-1].CreatedAt.Format(time.RFC3339Nano) + "|" + msgs[len(msgs)-1].Id.String()
	}

	return msgs, nextCursor, err
}
