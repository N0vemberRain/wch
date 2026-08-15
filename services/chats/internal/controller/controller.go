package controller

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"slices"
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

func (c *Controller) CreateGroupChat(
	ctx context.Context,
	chat *model.Chat,
	av *model.Avatar,
) (*model.Chat, error) {
	chat.ID = uuid.New()
	chat.CreatedAt = time.Now()
	chat.UpdatedAt = chat.CreatedAt

	err := c.repo.CreateChat(ctx, chat)
	if err != nil {
		return nil, err
	}

	id, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		log.Println("not ok")
		return nil, chats.ErrUserID
	}

	err = c.repo.AddParticipant(ctx, chat.ID, &model.ChatParticipant{
		ChatID: chat.ID,
		UserID: id,
		Role:   model.ChatParticipantOwner,
	})
	if err != nil {
		return nil, err
	}

	if av == nil || len(av.Data) == 0 {
		return chat, nil
	}

	av.OwnerID = chat.ID
	key, err := c.users.UpdateAvatarForChat(ctx, av)
	if err != nil {
		return nil, err
	}

	chat.AvatarKey = key
	err = c.repo.UpdateChat(ctx, chat)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (c *Controller) CreateDirectChat(
	ctx context.Context,
	chat *model.Chat,
	userID uuid.UUID,
) (*model.Chat, *model.Avatar, error) {
	// check if the user exists
	user, err := c.users.GetUserByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	chat.ID = uuid.New()
	chat.CreatedAt = time.Now()
	chat.UpdatedAt = chat.CreatedAt
	chat.AvatarKey = user.ID.String()

	err = c.repo.CreateChat(ctx, chat)
	if err != nil {
		return nil, nil, err
	}

	// add to participants who we are writing
	err = c.AddParticipant(ctx, chat.ID, &model.ChatParticipant{
		ChatID: chat.ID,
		UserID: userID,
		Role:   model.ChatParticipantAdmin,
	})
	if err != nil {
		return nil, nil, err
	}
	id, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		log.Printf("Current user: %v\n", id)
		return nil, nil, chats.ErrUserID
	}

	// add to articipants ourself
	err = c.AddParticipant(ctx, chat.ID, &model.ChatParticipant{
		ChatID: chat.ID,
		UserID: id,
		Role:   model.ChatParticipantAdmin,
	})
	if !ok {
		return nil, nil, chats.ErrUserID
	}

	// direct chat doesn't have name in db, it gets the name only when it returns to a user
	chat.Name = user.Username
	// get an avatar
	av, err := c.users.GetAvatarForDirectChatByKey(ctx, chat.AvatarKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Avatar for user %s not found. New chat will have no avatar", userID.String())
			return chat, nil, nil
		} else {
			log.Printf("Get Avatar Error: %v\n", err)
			return chat, nil, nil
		}
	}
	return chat, av, err
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

// Если чат - диалог, берем из БД его участников и
// присваиваем имя того участника, которые не является пользователем,
// например, зашли под игорем - имя чата маша, зашли под машей - имя чата игорь
func (c *Controller) getNameToDirectChat(ctx context.Context, chat *model.Chat) (
	string,
	error,
) {
	if chat.Type != model.ChatTypeDirect {
		panic("chat must be direct")
	}

	users, err := c.repo.GetParticipantsIDs(ctx, chat.ID)
	if err != nil {
		return "", err
	}

	log.Printf("Participants for chat %v: %v\n", chat.ID, users)
	if len(users) != 2 {
		panic("len(users) != 2")
	}

	loggedUser := ctx.Value("user_id").(uuid.UUID)
	log.Printf("LOGGED USER %s\n", loggedUser.String())
	if users[0] == ctx.Value("user_id").(uuid.UUID) {
		usr, err := c.users.GetUserByID(ctx, users[1])
		if err != nil {
			return "", err
		}

		return usr.Username, nil
	} else if users[1] == ctx.Value("user_id").(uuid.UUID) {
		usr, err := c.users.GetUserByID(ctx, users[0])
		if err != nil {
			return "", err
		}

		return usr.Username, nil
	} else {
		panic("one of direct chat participants must be a logged user")
	}
}

func (c *Controller) GetChatByID(ctx context.Context, chatID uuid.UUID) (
	*model.Chat,
	error,
) {
	if chatID == uuid.Nil {
		return nil, chats.ErrUserIDNil
	}

	chat, err := c.repo.GetChatByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if chat.Type == model.ChatTypeDirect {
		chat.Name, err = c.getNameToDirectChat(ctx, chat)
		if err != nil {
			return nil, err
		}
	}

	return chat, nil
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

func (c *Controller) ListParticipants(ctx context.Context, chatID uuid.UUID) ([]model.ChatParticipant, error) {
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

	avs, err := c.users.GetAvatarsForChats(ctx, ids)
	if err != nil {
		return nil, err
	}

	log.Printf("Found %d avatars\n", len(avs))

	participants := make([]model.ChatParticipant, len(users))

	if len(users) != len(ids) {
		panic("len(users) != len(ids)")
	}
	for i, u := range users {
		participants[i].UserID = u.ID
		participants[i].ChatID = chatID
		participants[i].Name = u.Username
		log.Printf("Participant name: %s\n", u.Username)
		// Now every one is admin
		participants[i].Role = model.ChatParticipantAdmin

		for j := range avs {
			log.Printf("Avatar for owner %v\n", avs[j].OwnerID)
			if avs[j].OwnerID == u.ID {
				log.Printf("Avatar found for owner %v\n", avs[j].OwnerID)
				participants[i].Avatar = avs[j]
			}
		}
	}

	return participants, nil
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

	chats, err := c.repo.GetChatsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i, chat := range chats {
		if chat.Type != model.ChatTypeDirect {
			continue
		}

		chats[i].Name, err = c.getNameToDirectChat(ctx, &chat)
		if err != nil {
			return nil, err
		}
	}

	return chats, nil
}

func (c *Controller) ListAvatarsForGroupChats(ctx context.Context, ids []uuid.UUID) (
	[]model.Avatar,
	error,
) {
	if len(ids) == 0 {
		return nil, chats.ErrChatIDNil
	}

	if slices.Contains(ids, uuid.Nil) {
		return nil, chats.ErrChatIDNil
	}

	return c.users.GetAvatarsForChats(ctx, ids)
}

func (c *Controller) ListAvatarsForDirectChats(ctx context.Context, ids []uuid.UUID) (
	[]model.Avatar,
	error,
) {
	if len(ids) == 0 {
		return nil, chats.ErrChatIDNil
	}

	if slices.Contains(ids, uuid.Nil) {
		return nil, chats.ErrChatIDNil
	}

	// Получаем id пользователей,
	// с которыми имеется переписка у текущего пользователя
	recievAndChat := make(map[uuid.UUID]uuid.UUID)
	var recieversIDs []uuid.UUID
	for _, chatID := range ids {
		id, err := c.getRecieverIDForDirectChat(ctx, chatID)
		if err != nil {
			return nil, err
		}

		recieversIDs = append(recieversIDs, id)
		recievAndChat[id] = chatID
	}

	// возвращаем аватары для пользователей-приемников сообщений
	avs, err := c.users.GetAvatarsForChats(ctx, recieversIDs)
	if err != nil {
		return nil, err
	}
	// заменяем owner id, который был равен user_id, на chat_id
	for i := range avs {
		avs[i].OwnerID = recievAndChat[avs[i].OwnerID]
	}

	return avs, nil
}

// Получаем id пользователей,
// с которыми имеется переписка у текущего пользователя
func (c *Controller) getRecieverIDForDirectChat(ctx context.Context, chatID uuid.UUID) (
	uuid.UUID,
	error,
) {
	users, err := c.repo.GetParticipantsIDs(ctx, chatID)
	if err != nil {
		return uuid.Nil, err
	}

	log.Printf("Participants for chat %v: %v\n", chatID, users)
	if len(users) != 2 {
		panic(len(users) != 2)
	}

	loggedUserID := ctx.Value("user_id").(uuid.UUID)
	log.Printf("LOGGED USER %s\n", loggedUserID.String())
	var recieverID uuid.UUID
	if users[0] == ctx.Value("user_id").(uuid.UUID) {
		recieverID = users[1]
	} else if users[1] == ctx.Value("user_id").(uuid.UUID) {
		recieverID = users[0]
	} else {
		panic("one of direct chat participants must be logged user")
	}

	return recieverID, nil
}
