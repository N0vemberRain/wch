package grpc

import (
	"context"
	"time"

	chatspb "wch/gen/chats/v1"
	"wch/services/chats/internal/controller"
	chats "wch/services/chats/internal/domain"
	"wch/services/chats/internal/domain/model"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrRequestIsEmpty = status.Error(codes.InvalidArgument, "request is empty")

type Handler struct {
	ctrl *controller.Controller
	chatspb.ChatsServiceServer
}

func NewHandler(ctrl *controller.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

func (h *Handler) CreateChat(ctx context.Context, req *chatspb.CreateChatRequest) (
	*chatspb.ChatResponse,
	error,
) {
	if req == nil {
		return nil, ErrRequestIsEmpty
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

	return &chatspb.ChatResponse{}, nil
}

func (h *Handler) GetChatByID(ctx context.Context, req *chatspb.GetChatByIDRequest) (
	*chatspb.GetChatResponse,
	error,
) {
	if req == nil {
		return nil, ErrRequestIsEmpty
	}

	id, err := uuid.Parse(req.ChatId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	chat, err := h.ctrl.GetChatByID(ctx, id)
	if err != nil {
		if err == chats.ErrChatNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &chatspb.GetChatResponse{
		Chat: model.ChatToProto(chat),
	}, nil
}

func (h *Handler) GetChatByName(ctx context.Context, req *chatspb.GetChatByNameRequest) (
	*chatspb.GetChatResponse,
	error,
) {
	if req == nil {
		return nil, ErrRequestIsEmpty
	}

	if req.Name == "" {
		return nil, chats.ErrChatNameIsEmpty
	}
	chat, err := h.ctrl.GetChatByName(ctx, req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &chatspb.GetChatResponse{
		Chat: model.ChatToProto(chat),
	}, nil
}

func (h *Handler) UpdateChat(ctx context.Context, req *chatspb.UpdateChatRequest) (
	*chatspb.ChatResponse,
	error,
) {
	if req == nil {
		return nil, ErrRequestIsEmpty
	}

	chat, err := model.ChatFromProto(req.Chat)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = h.ctrl.UpdateChat(ctx, chat)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &chatspb.ChatResponse{}, nil
}

func (h *Handler) DeleteChat(ctx context.Context, req *chatspb.DeleteChatRequest) (
	*chatspb.ChatResponse,
	error,
) {
	if req == nil {
		return nil, ErrRequestIsEmpty
	}

	id, err := uuid.Parse(req.ChatId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = h.ctrl.DeleteChat(ctx, id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &chatspb.ChatResponse{}, nil
}

func (h *Handler) AddParticipant(ctx context.Context, req *chatspb.AddParticipantRequest) (
	*chatspb.ParticipantResponse,
	error,
) {
	if req == nil {
		return nil, ErrRequestIsEmpty
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
	err = h.ctrl.AddParticipant(ctx, chat_id, participant)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &chatspb.ParticipantResponse{Ok: true}, nil
}

func (h *Handler) RemoveParticipant(ctx context.Context, req *chatspb.RemoveParticipantRequest) (
	*chatspb.ParticipantResponse,
	error,
) {
	if req == nil {
		return nil, ErrRequestIsEmpty
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

	return &chatspb.ParticipantResponse{Ok: true}, nil
}

func (h *Handler) ListParticipants(ctx context.Context, req *chatspb.ListParticipantsRequest) (
	*chatspb.ListParticipantsResponse,
	error,
) {
	if req == nil {
		return nil, ErrRequestIsEmpty
	}

	chatId, err := uuid.Parse(req.ChatId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	participants, err := h.ctrl.ListParticipants(ctx, chatId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp := &chatspb.ListParticipantsResponse{}

	for _, p := range participants {
		resp.Users = append(resp.Users, model.UserToProto(&p))
	}

	return resp, nil
}

func (h *Handler) CheckParticipant(ctx context.Context, req *chatspb.CheckParticipantRequest) (
	*chatspb.CheckParticipantResponse,
	error,
) {
	if req == nil {
		return nil, ErrRequestIsEmpty
	}

	chatID, err := uuid.Parse(req.ChatId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	participant, err := h.ctrl.GetParticipant(ctx, chatID, userID)
	if err != nil && err == chats.ErrParticipateNotFound {
		return &chatspb.CheckParticipantResponse{
			IsParticipant: false,
			Role:          model.ChatParticipantUnknown,
		}, nil
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &chatspb.CheckParticipantResponse{
		IsParticipant: true,
		Role:          chatspb.ChatParticipantRole(participant.Role),
	}, nil
}
