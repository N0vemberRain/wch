package model

import (
	msgspb "wch/gen/msgs/v1"
	"wch/services/msgs/internal/domain"

	"github.com/google/uuid"
	timeconv "google.golang.org/protobuf/types/known/timestamppb"
)

func MessageToProto(m *Message) *msgspb.Message {
	return &msgspb.Message{
		Id:        m.Id.String(),
		ChatId:    m.ChatId.String(),
		SenderId:  m.SenderId.String(),
		Content:   m.Content,
		CreatedAt: timeconv.New(m.CreatedAt),
	}
}

func MessageFromProto(m *msgspb.Message) (*Message, error) {
	id, err := uuid.Parse(m.Id)
	if err != nil {
		return nil, domain.ErrMessageID
	}
	chat_id, err := uuid.Parse(m.ChatId)
	if err != nil {
		return nil, domain.ErrChatID
	}
	sender_id, err := uuid.Parse(m.SenderId)
	if err != nil {
		return nil, domain.ErrSenderID
	}
	return &Message{
		Id:        id,
		ChatId:    chat_id,
		SenderId:  sender_id,
		Content:   m.Content,
		CreatedAt: m.CreatedAt.AsTime(),
	}, nil
}
