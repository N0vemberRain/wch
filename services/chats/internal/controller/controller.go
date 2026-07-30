package controller

import (
	"context"
	"time"

	chats "wch/services/chats/internal/domain"
	"wch/services/chats/internal/domain/model"
	"wch/services/chats/internal/domain/ports"

	"github.com/google/uuid"
)

type Controller struct {
	repo  ports.Repository
	users ports.UserProvider
}

func NewChatController(repo ports.Repository, users ports.UserProvider) *Controller {
	return &Controller{
		repo:  repo,
		users: users,
	}
}

func (c *Controller) CreateChat(ctx context.Context, chat *model.Chat) error {
	chat.ID = uuid.New()
	chat.CreatedAt = time.Now()
	chat.UpdatedAt = chat.CreatedAt

	return c.repo.CreateChat(ctx, chat)
}

func (c *Controller) UpdateChat(ctx context.Context, chat *model.Chat, av_bytes []byte) error {
	chat.UpdatedAt = time.Now()

	if chat.Type == model.ChatTypeGroup {
		if len(av_bytes) != 0 {
			key, err := c.users.UpdateAvatarForChat(
				ctx,
				&model.Avatar{
					OwnerID:  chat.ID,
					Data:     av_bytes,
					MimeType: "PNG",
				},
			)

			if err != nil {
				return err
			}

			chat.AvatarKey = key
		}
	}

	return c.repo.UpdateChat(ctx, chat)
}

func (c *Controller) GetChatByID(ctx context.Context, chatID uuid.UUID) (*model.Chat, error) {
	if chatID == uuid.Nil {
		return nil, chats.ErrUserIDNil
	}

	return c.repo.GetChatByID(ctx, chatID)
}

func (c *Controller) GetChatByName(ctx context.Context, name string) (*model.Chat, error) {
	return c.repo.GetChatByName(ctx, name)
}
func (c *Controller) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	if chatID == uuid.Nil {
		return chats.ErrUserIDNil
	}

	return c.repo.DeleteChat(ctx, chatID)
}

func (c *Controller) AddParticipant(ctx context.Context, chatID uuid.UUID, p *model.ChatParticipant) error {
	if p.ChatID == uuid.Nil {
		return chats.ErrChatIDNil
	}
	if p.UserID == uuid.Nil {
		return chats.ErrUserIDNil
	}

	p.JoinedAt = time.Now()

	return c.repo.AddParticipant(ctx, chatID, p)
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

func (c *Controller) ListParticipants(ctx context.Context, chatID uuid.UUID) ([]model.User, error) {
	if chatID == uuid.Nil {
		return nil, chats.ErrChatIDNil
	}

	ids, err := c.repo.GetParticipantsIDs(ctx, chatID)
	if err != nil {
		return nil, err
	}

	users, err := c.users.GetUsersByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (c *Controller) GetParticipant(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (
	*model.ChatParticipant,
	error,
) {
	if chatID == uuid.Nil {
		return nil, chats.ErrChatIDNil
	}

	if userID == uuid.Nil {
		return nil, chats.ErrUserIDNil
	}

	p, err := c.repo.GetParticipant(ctx, chatID, userID)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (c *Controller) ListChatsForUser(ctx context.Context, userID uuid.UUID) (
	[]model.Chat,
	error,
) {
	if userID == uuid.Nil {
		return nil, chats.ErrUserIDNil
	}

	return c.repo.GetChatsForUser(ctx, userID)
}

func (c *Controller) ListAvatarsForChats(ctx context.Context, ids []uuid.UUID) (
	[]model.Avatar,
	error,
) {
	if len(ids) == 0 {
		return nil, chats.ErrChatIDNil
	}

	for _, id := range ids {
		if id == uuid.Nil {
			return nil, chats.ErrChatIDNil
		}
	}

	return c.users.GetAvatarsForChats(ctx, ids)
}
