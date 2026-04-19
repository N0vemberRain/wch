package ports

import (
	"context"
	"wch/services/nots/internal/domain"
)

type NotificationSender interface {
	Send(context.Context, domain.Notification) error
}
