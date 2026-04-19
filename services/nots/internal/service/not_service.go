package service

import (
	"context"
	"wch/pkg/events"
	"wch/services/nots/internal/domain"
	"wch/services/nots/internal/ports"
)

type NotificationService interface {
	Notify(context.Context, events.MessageSentEvent) error
}

type notificationService struct {
	sender       ports.NotificationSender
	chatProvider ports.ChatProvider
}

func NewNotificationService(sender ports.NotificationSender, chatProvider ports.ChatProvider) NotificationService {
	return &notificationService{sender: sender, chatProvider: chatProvider}
}

func (ns *notificationService) Notify(ctx context.Context, event events.MessageSentEvent) error {
	ids, err := ns.chatProvider.GetUserIDsInChat(ctx, event.ChatID)
	if err != nil {
		return err
	}

	for _, id := range ids {
		not := domain.Notification{
			UserID:  id,
			ChatID:  event.ChatID,
			Content: event.Content,
		}

		ns.sender.Send(ctx, not)
	}

	return nil
}
