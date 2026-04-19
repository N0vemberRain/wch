package adapters

import (
	"context"
	notspb "wch/gen/nots/v1"
	"wch/services/nots/internal/domain"
	"wch/services/nots/internal/streams"
)

type GRPCSender struct {
	hub *streams.StreamHub
}

func NewGRPCSender(hub *streams.StreamHub) *GRPCSender {
	return &GRPCSender{
		hub: hub,
	}
}

func (s *GRPCSender) Send(ctx context.Context, n domain.Notification) error {
	grpcNot := &notspb.Notification{
		UserId:  n.UserID.String(),
		ChatId:  n.ChatID.String(),
		Content: n.Content,
	}

	return s.hub.SendToUser(n.UserID.String(), grpcNot)
}
