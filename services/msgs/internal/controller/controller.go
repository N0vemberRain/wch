package controller

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"wch/pkg/events"
	"wch/services/msgs/internal/domain"
	"wch/services/msgs/internal/domain/model"
	"wch/services/msgs/internal/domain/ports"
	"wch/services/msgs/internal/domain/shared"

	"github.com/google/uuid"
)

type Controller struct {
	repo        ports.MessageRepository
	permChecker *shared.PermissionChecker
	publisher   ports.EventPublisher
}

func NewMessageController(
	repo ports.MessageRepository,
	permChecker *shared.PermissionChecker,
	publisher ports.EventPublisher,
) *Controller {
	return &Controller{
		repo:        repo,
		permChecker: permChecker,
		publisher:   publisher,
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
	c.repo.Save(ctx, msg)
	event := &events.MessageSentEvent{
		Type:      "MESSAGE_SENT",
		ID:        msg.Id,
		ChatID:    msg.ChatId,
		SenderID:  msg.SenderId,
		Content:   msg.Content,
		Timestamp: msg.CreatedAt,
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return c.publisher.Publish(ctx, "message_sent", msg.ChatId.String(), payload)
}

func (c *Controller) List(ctx context.Context, chatID uuid.UUID, limit int, cursor *model.Cursor) ([]model.Message, string, error) {
	if chatID == uuid.Nil {
		return nil, "", domain.ErrChatIsEmpty
	}

	log.Println("msgs.Controller.List: checking if chat exists")
	ok, err := c.permChecker.DoesChatExists(ctx, chatID)
	if err != nil {
		return nil, "", err
	}
	if !ok {
		return nil, "", domain.ErrChatNotFound
	}

	log.Println("msgs.Controller.List: calling Repository.List")
	msgs, err := c.repo.List(ctx, chatID, limit, cursor)
	var nextCursor string
	if err == nil && len(msgs) >= limit {
		nextCursor = msgs[len(msgs)-1].CreatedAt.Format(time.RFC3339Nano) + "|" + msgs[len(msgs)-1].Id.String()
	}

	return msgs, nextCursor, err
}
