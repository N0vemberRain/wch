package adapters

import (
	"context"
	"log"
	"wch/pkg/events"
	"wch/services/nots/internal/service"
)

type MsgsHandlerGRPC struct {
	service service.NotificationService
}

func NewMsgsHandlerGRPC(service service.NotificationService) *MsgsHandlerGRPC {
	return &MsgsHandlerGRPC{
		service: service,
	}
}

func (s *MsgsHandlerGRPC) Handle(ctx context.Context, e events.MessageSentEvent) error {
	log.Println("Handling kafka event: ", e)

	return s.service.Notify(ctx, e)
}
