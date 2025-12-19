package grpc

import (
	"context"
	"time"

	"wch/gen"
	"wch/services/chats"
	"wch/services/chats/internal/controller"
	"wch/services/chats/pkg/model"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	ctrl *controller.Controller
}

func (h *Handler) CreateChat(ctx context.Context, req *gen.CreateChatRequest) (
	*gen.ChatResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, chats.ErrRequestIsEmpty.Error())
	}

	t, err := model.ChatTypeFromString(req.Chat.Type)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, chats.ErrChatType.Error())
	}
	if t != model.ChatTypeDirect && req.Chat.Name == "" {
		return nil, status.Error(codes.InvalidArgument, chats.ErrChatNameIsEmpty.Error())
	}

	chat := model.NewChat(req.Chat.Name, t)
	err = h.ctrl.CreateChat(ctx, chat)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.ChatResponse{}, nil
}

func (h *Handler) GetChat(ctx context.Context, req *gen.GetChatRequest) (
	*gen.GetChatResponse,
	error,
) {
	if req == nil {
		return nil, chats.ErrRequestIsEmpty
	}

	id, err := uuid.Parse(req.ChatId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	chat, err := h.ctrl.GetChat(ctx, id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &gen.GetChatResponse{
		Chat: model.ChatToProto(chat),
	}, nil
}

func (h *Handler) UpdateChat(ctx context.Context, req *gen.UpdateChatRequest) (
	*gen.ChatResponse,
	error,
) {
	if req == nil {
		return nil, chats.ErrRequestIsEmpty
	}

	chat, err := model.ChatFromProto(req.Chat)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = h.ctrl.UpdateChat(ctx, chat)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &gen.ChatResponse{}, nil
}

func (h *Handler) DeleteChat(ctx context.Context, req *gen.DeleteChatRequest) (
	*gen.ChatResponse,
	error,
) {
	if req == nil {
		return nil, chats.ErrRequestIsEmpty
	}

	id, err := uuid.Parse(req.ChatId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = h.ctrl.DeleteChat(ctx, id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &gen.ChatResponse{}, nil
}

func (h *Handler) AddParticipant(ctx context.Context, req *gen.AddParticipantRequest) (
	*gen.ParticipantResponse,
	error,
) {
	if req == nil {
		return nil, chats.ErrRequestIsEmpty
	}

	if req.ChatId == "" {
		return nil, status.Error(codes.InvalidArgument, chats.ErrChatID.Error())
	} else if req.Participant.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, chats.ErrUserID.Error())
	} else if req.Participant.Role == model.ChatParticipantUnknown {
		return nil, status.Error(codes.InvalidArgument, "participant's role is undefined")
	}

	chat_id, err := uuid.Parse(req.ChatId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	user_id, err := uuid.Parse(req.Participant.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	participant, err := model.NewChatParticipant(
		chat_id,
		user_id,
		model.ChatParticipantRole(req.Participant.Role),
		time.Now(),
	)
	err = h.ctrl.AddParticipant(ctx, participant)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &gen.ParticipantResponse{Ok: true}, nil
}

func (h *Handler) RemoveParticipant(ctx context.Context, req *gen.RemoveParticipantRequest) (
	*gen.ParticipantResponse,
	error,
) {
	if req == nil {
		return nil, chats.ErrRequestIsEmpty
	}

	chatId, err := uuid.Parse(req.ChatId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = h.ctrl.RemoveParticipant(ctx, chatId, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.ParticipantResponse{Ok: true}, nil
}

func (h *Handler) ListParticipants(ctx context.Context, req *gen.ListParticipantsRequest) (
	*gen.ListParticipantsResponse,
	error,
) {
	if req == nil {
		return nil, chats.ErrRequestIsEmpty
	}

	chatId, err := uuid.Parse(req.ChatId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	participants, err := h.ctrl.ListParticipants(ctx, chatId)

	resp := &gen.ListParticipantsResponse{}

	for _, p := range participants {
		resp = append(resp, p)
	}
}
