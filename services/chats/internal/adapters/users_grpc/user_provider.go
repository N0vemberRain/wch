package usersgrpc

import (
	"context"

	userspb "wch/gen/users/v1"
	"wch/services/chats/internal/domain/model"

	"github.com/google/uuid"
)

type UserProviderGRPC struct {
	client userspb.UsersServiceClient
}

func (up *UserProviderGRPC) GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]model.User, error) {
	req := &userspb.GetUsersByIDsRequest{
		Id: model.IDsToString(ids),
	}

	resp, err := up.client.GetUsersByIDs(ctx, req)
	if err != nil {
		return nil, err
	}

	users := make([]model.User, len(resp.Users))
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
