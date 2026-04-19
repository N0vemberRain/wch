package chatsgrpc

import (
	"context"

	chatspb "wch/gen/chats/v1"

	"github.com/google/uuid"
)

type ChatProviderGRPC struct {
	client chatspb.ChatsServiceClient
}

func NewChatProvider(client chatspb.ChatsServiceClient) *ChatProviderGRPC {
	return &ChatProviderGRPC{
		client: client,
	}
}

func (cp *ChatProviderGRPC) GetUserIDsInChat(ctx context.Context, chatID uuid.UUID) (
	uuid.UUIDs,
	error,
) {
	resp, err := cp.client.ListParticipants(ctx, &chatspb.ListParticipantsRequest{
		ChatId: chatID.String(),
	})

	var ids uuid.UUIDs
	if err != nil {
		return ids, err
	}

	for _, user := range resp.Users {
		ids = append(ids, uuid.MustParse(user.Id))
	}

	return ids, nil
}
