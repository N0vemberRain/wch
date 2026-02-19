package chatsgrpc

import (
	"context"

	chatspb "wch/gen/chats/v1"
	"wch/services/msgs/internal/domain/model"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatProviderGRPC struct {
	client chatspb.ChatsServiceClient
}

func NewChatProvider(client chatspb.ChatsServiceClient) *ChatProviderGRPC {
	return &ChatProviderGRPC{
		client: client,
	}
}

func (cp *ChatProviderGRPC) IsUserParticipant(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (
	bool, model.ChatParticipantRole, error,
) {
	resp, err := cp.client.CheckParticipant(ctx, &chatspb.CheckParticipantRequest{
		ChatId: chatID.String(),
		UserId: userID.String(),
	})
	if err != nil {
		return false, model.ChatParticipantUnknown, err
	}

	if !resp.IsParticipant {
		return false, model.ChatParticipantUnknown, nil
	}

	return true, model.ChatParticipantRoleFromProto(resp.Role), nil
}

func (cp *ChatProviderGRPC) ChatExists(ctx context.Context, chatID uuid.UUID) (bool, error) {
	_, err := cp.client.GetChatByID(ctx, &chatspb.GetChatByIDRequest{
		ChatId: chatID.String(),
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}
