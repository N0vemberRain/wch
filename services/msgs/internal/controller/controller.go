package controller

import (
	"context"
	"time"
	"wch/services/msgs/internal/domain"
	"wch/services/msgs/internal/domain/model"
	"wch/services/msgs/internal/domain/ports"
	"wch/services/msgs/internal/domain/shared"

	"github.com/google/uuid"
)

type Controller struct {
	repo        ports.MessageRepository
	permChecker *shared.PermissionChecker
}

func NewMessageController(repo ports.MessageRepository, permChecker *shared.PermissionChecker) *Controller {
	return &Controller{
		repo:        repo,
		permChecker: permChecker,
	}
}

func (c *Controller) Save(ctx context.Context, msg *model.Message) error {
	canSend, err := c.permChecker.CanSendMessage(ctx, msg.ChatId, msg.SenderId)
	if err != nil {
		return err
	}
	if !canSend {
		return domain.ErrNotAllowedToSend
	}
	msg.Id = uuid.New()
	msg.CreatedAt = time.Now()
	return c.repo.Save(ctx, msg)
}

func (c *Controller) List(ctx context.Context, chatID uuid.UUID, limit int, cursor *model.Cursor) ([]model.Message, string, error) {
	if chatID == uuid.Nil {
		return nil, "", domain.ErrChatIsEmpty
	}

	ok, err := c.permChecker.DoesChatExists(ctx, chatID)
	if err != nil {
		return nil, "", err
	}
	if !ok {
		return nil, "", domain.ErrChatNotFound
	}

	msgs, err := c.repo.List(ctx, chatID, limit, cursor)
	var nextCursor string
	if err == nil && len(msgs) >= limit {
		nextCursor = msgs[len(msgs)-1].CreatedAt.Format(time.RFC3339Nano) + "|" + msgs[len(msgs)-1].Id.String()
	}

	return msgs, nextCursor, err
}
