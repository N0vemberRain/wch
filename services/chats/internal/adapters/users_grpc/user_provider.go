package usersgrpc

import (
	"context"
	"log"

	userspb "wch/gen/users/v1"
	"wch/pkg/auth"
	"wch/services/chats/internal/domain/model"

	"github.com/google/uuid"
)

type UserProviderGRPC struct {
	client userspb.UsersServiceClient
}

func NewUserProvider(client userspb.UsersServiceClient) *UserProviderGRPC {
	return &UserProviderGRPC{
		client: client,
	}
}

func (up *UserProviderGRPC) GetUserByID(ctx context.Context, id uuid.UUID) (
	*model.User,
	error,
) {
	req := &userspb.GetUserByIDRequest{
		UserId: id.String(),
	}

	outCtx := auth.ForwardAuthContext(ctx)
	resp, err := up.client.GetUserByID(outCtx, req)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:        uuid.MustParse(resp.User.Id),
		Username:  resp.User.Username,
		AvatarURL: resp.User.AvatarUrl,
		Status:    resp.User.Status,
	}, nil
}

func (up *UserProviderGRPC) GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]model.User, error) {
	req := &userspb.GetUsersByIDsRequest{
		Id: model.IDsToString(ids),
	}

	outCtx := auth.ForwardAuthContext(ctx)
	log.Printf("GetUsersByIDs: %v\n", ctx.Value("user_id"))
	resp, err := up.client.GetUsersByIDs(outCtx, req)
	if err != nil {
		return nil, err
	}

	users := make([]model.User, 0)
	for _, u := range resp.Users {
		u := model.User{
			ID:        uuid.MustParse(u.Id),
			Username:  u.Username,
			AvatarURL: u.AvatarUrl,
		}
		users = append(users, u)
	}

	return users, nil
}

func (up *UserProviderGRPC) GetAvatarForDirectChatByKey(ctx context.Context, key string) (
	*model.Avatar,
	error,
) {
	// in direct chats key is a user'd id
	id, err := uuid.Parse(key)
	if err != nil {
		return nil, err
	}
	req := &userspb.GetAvatarRequest{
		UserId: id.String(),
	}
	outCtx := auth.ForwardAuthContext(ctx)
	resp, err := up.client.GetAvatarForUser(outCtx, req)
	if err != nil {
		return nil, err
	}

	// ownerId, err := uuid.Parse(resp.Avatar.OwnerId)
	// if err != nil {
	// 	return nil, err
	// }
	return &model.Avatar{
		Data:     resp.Avatar.Data,
		MimeType: resp.Avatar.MimeType,
		OwnerID:  id,
	}, nil
}

func (up *UserProviderGRPC) GetAvatarsForChats(ctx context.Context, ids []uuid.UUID) (
	[]model.Avatar,
	error,
) {
	req := &userspb.GetAvatarsForChatsRequest{
		Ids: model.IDsToString(ids),
	}

	outCtx := auth.ForwardAuthContext(ctx)
	log.Printf("GetAvatarsForChats: %v\n", ctx.Value("user_id"))
	resp, err := up.client.GetAvatarsForChats(outCtx, req)
	if err != nil {
		return nil, err
	}

	avatars := make([]model.Avatar, 0)
	for i, aProto := range resp.Avatars {
		a := model.Avatar{
			Data:     aProto.Data,
			MimeType: aProto.MimeType,
			OwnerID:  uuid.MustParse(aProto.OwnerId),
		}

		avatars = append(avatars, a)
		log.Printf("Avatar %d for chat %s\n", i, a.OwnerID)
	}

	return avatars, nil
}

func (up *UserProviderGRPC) UpdateAvatarForChat(
	ctx context.Context,
	av *model.Avatar,
) (string, error) {
	req := &userspb.UpdateAvatarForChatRequest{
		ChatId: av.OwnerID.String(),
		Avatar: av.Data,
	}

	outCtx := auth.ForwardAuthContext(ctx)

	resp, err := up.client.UpdateAvatarForChat(outCtx, req)
	if err != nil {
		return "", err
	}

	return resp.Key, nil
}
