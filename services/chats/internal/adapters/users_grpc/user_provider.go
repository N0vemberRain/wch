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

func (up *UserProviderGRPC) GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]model.User, error) {
	req := &userspb.GetUsersByIDsRequest{
		Id: model.IDsToString(ids),
	}

	log.Printf("UpdateAvatarForChat: %v\n", ctx.Value("user_id"))
	resp, err := up.client.GetUsersByIDs(ctx, req)
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

func (up *UserProviderGRPC) GetAvatarsForChats(ctx context.Context, ids []uuid.UUID) (
	[]model.Avatar,
	error,
) {
	req := &userspb.GetAvatarsForChatsRequest{
		Ids: model.IDsToString(ids),
	}

	log.Printf("UpdateAvatarForChat: %v\n", ctx.Value("user_id"))
	resp, err := up.client.GetAvatarsForChats(ctx, req)
	if err != nil {
		return nil, err
	}

	avatars := make([]model.Avatar, 0)
	for _, aProto := range resp.Avatars {
		a := model.Avatar{
			Data:     aProto.Data,
			MimeType: aProto.MimeType,
			OwnerID:  uuid.MustParse(aProto.OwnerId),
		}

		avatars = append(avatars, a)
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
