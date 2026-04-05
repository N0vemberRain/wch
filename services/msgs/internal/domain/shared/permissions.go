package shared

import (
	"context"
	"errors"
	"log"
	"wch/services/msgs/internal/domain/ports"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
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
	log.Println("PermissionChecker.CanSendMessage: entering...")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false, errors.New("PermissionChecker.DoesChatExists: metadata doesn't exists")
	}
	ok, _, err := pc.chatProvider.IsUserParticipant(metadata.NewOutgoingContext(ctx, md), chatID, userID)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	log.Println("PermissionChecker.CanSendMessage: exiting...")
	return true, nil
}

func (pc *PermissionChecker) DoesChatExists(ctx context.Context, chatID uuid.UUID) (bool, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false, errors.New("PermissionChecker.DoesChatExists: metadata doesn't exists")
	}

	return pc.chatProvider.ChatExists(metadata.NewOutgoingContext(ctx, md), chatID)
}
