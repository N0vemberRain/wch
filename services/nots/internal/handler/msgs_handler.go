package handler

import (
	"context"
	"wch/pkg/events"
)

type MsgsHandler interface {
	Handle(ctx context.Context, event events.MessageSentEvent) error
}
