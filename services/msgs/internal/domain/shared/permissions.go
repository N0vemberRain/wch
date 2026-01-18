package shared

import (
	"context"
	"wch/services/msgs/internal/domain/ports"

	"github.com/google/uuid"
)

type PermissionChecker struct {
	chatProvider ports.ChatProvider
}

func NewPermissionChecker(chatProvider ports.ChatProvider) *PermissionChecker {
	return &PermissionChecker{
		chatProvider: chatProvider,
	}
}

func (pc *PermissionChecker) CanSendMessage(
	ctx context.Context, chatID uuid.UUID, userID uuid.UUID,
) (bool, error) {
	ok, _, err := pc.chatProvider.IsUserParticipant(ctx, chatID, userID)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}
