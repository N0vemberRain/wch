package grpc

import (
	"context"
	"errors"
	msgspb "wch/gen/msgs/v1"
	"wch/services/msgs/internal/controller"
	"wch/services/msgs/internal/domain"
	"wch/services/msgs/internal/domain/model"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrRequestIsEmpty = errors.New("request is empty")

type Handler struct {
	msgspb.MessagesServiceServer
	ctrl *controller.Controller
}

func NewHandler(ctrl *controller.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

func (h *Handler) SendMessage(ctx context.Context, req *msgspb.SendMessageRequest) (
	*msgspb.MessageResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, ErrRequestIsEmpty.Error())
	}

	if req.ChatId == "" {
		return nil, status.Error(codes.InvalidArgument, domain.ErrChatIsEmpty.Error())
	}

	if req.SenderId == "" {
		return nil, status.Error(codes.InvalidArgument, domain.ErrSenderIDIsEmpty.Error())
	}

	chat_id, err := uuid.Parse(req.ChatId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, domain.ErrChatID.Error())
	}
	sender_id, err := uuid.Parse(req.SenderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, domain.ErrSenderID.Error())
	}
	msg := &model.Message{
		ChatId:   chat_id,
		SenderId: sender_id,
		Content:  req.Content,
	}

	err = h.ctrl.Save(ctx, msg)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &msgspb.MessageResponse{}, nil
}

func (h *Handler) ListMessages(ctx context.Context, req *msgspb.ListMessagesRequest) (
	*msgspb.ListMessagesResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, ErrRequestIsEmpty.Error())
	}

	if req.ChatId == "" {
		return nil, status.Error(codes.InvalidArgument, domain.ErrChatIsEmpty.Error())
	}

	chat_id, err := uuid.Parse(req.ChatId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, domain.ErrChatID.Error())
	}

	if req.Limit == 0 {
		return nil, status.Error(codes.InvalidArgument, "limit is undefined")
	}

	msgs := make([]model.Message, 0)
	var nextCursor string
	if req.Cursor == "" {
		msgs, nextCursor, err = h.ctrl.List(ctx, chat_id, int(req.Limit), nil)
		if err != nil {
			if err == domain.ErrChatNotFound {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else {
		cursor, err := model.NewCursor(req.Cursor)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		msgs, nextCursor, err = h.ctrl.List(ctx, chat_id, int(req.Limit), cursor)
		if err != nil {
			if err == domain.ErrChatNotFound {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	resp := &msgspb.ListMessagesResponse{}
	for _, msg := range msgs {
		resp.Msgs = append(resp.Msgs, model.MessageToProto(&msg))
	}
	resp.NextCursor = nextCursor

	return resp, nil
}
