package controller

import (
	"context"
	"time"

	"wch/services/chats"
	"wch/services/chats/internal/repository"
	"wch/services/chats/pkg/model"

	users "wch/services/users/pkg/model"

	"github.com/google/uuid"
)

type Controller struct {
	repo repository.Repository
}

func (c *Controller) CreateChat(ctx context.Context, chat *model.Chat) error {
	chat.ID = uuid.New()
	chat.CreatedAt = time.Now()
	chat.UpdatedAt = chat.CreatedAt

	return c.repo.CreateChat(ctx, chat)
}

func (c *Controller) UpdateChat(ctx context.Context, chat *model.Chat) error {
	chat.UpdatedAt = time.Now()

	return c.repo.UpdateChat(ctx, chat)
}

func (c *Controller) GetChat(ctx context.Context, chatID uuid.UUID) (*model.Chat, error) {
	if chatID == uuid.Nil {
		return nil, chats.ErrUserIDNil
	}

	return c.repo.GetChat(ctx, chatID)
}

func (c *Controller) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	if chatID == uuid.Nil {
		return chats.ErrUserIDNil
	}

	return c.repo.DeleteChat(ctx, chatID)
}

func (c *Controller) AddParticipant(ctx context.Context, p *model.ChatParticipant) error {
	if p.ChatID == uuid.Nil {
		return chats.ErrChatIDNil
	}
	if p.UserID == uuid.Nil {
		return chats.ErrUserIDNil
	}

	p.JoinedAt = time.Now()

	return c.repo.AddParticipant(ctx, p)
}

func (c *Controller) RemoveParticipant(ctx context.Context, chatID, userID uuid.UUID) error {
	if chatID == uuid.Nil {
		return chats.ErrChatIDNil
	}
	if userID == uuid.Nil {
		return chats.ErrUserIDNil
	}

	return c.repo.RemoveParticipant(ctx, chatID, userID)
}

func (c *Controller) ListParticipants(ctx context.Context, chatID uuid.UUID) ([]*users.User, error) {
	if chatID == uuid.Nil {
		return nil, chats.ErrChatIDNil
	}

	ids, err := c.repo.GetParticipantsIDs(ctx, chatID)
	if err != nil {
		return nil, err
	}

	ids_str := make([]string, len(ids))
	for _, id := range ids {
		ids_str = append(ids_str, id.String())
	}
	resp, err := c.usersClient.GetUsersByIDs(ctx, &gen.GetUsersByIDsRequest{
		Id: ids_str,
	})

	if err != nil {
		return nil, err
	}

	return resp.Users, nil
}
