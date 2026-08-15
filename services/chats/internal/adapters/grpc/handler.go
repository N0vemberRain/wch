package grpc

import (
	"context"
	"log"
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

func (h *Handler) CreateGroupChat(ctx context.Context, req *chatspb.CreateGroupChatRequest) (
	*chatspb.ChatResponse,
	error,
) {
	if req == nil {
		return nil, ErrRequestIsEmpty
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, chats.ErrChatNameIsEmpty.Error())
	}

	chat := model.NewChat(req.Name, model.ChatTypeGroup)
	log.Printf("Chat: %v\t", chat)
	av := &model.Avatar{}
	if req.Avatar != nil {
		av.Data = req.Avatar.Data
		av.MimeType = req.Avatar.MimeType
	}
	chat, err := h.ctrl.CreateGroupChat(ctx, chat, av)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &chatspb.ChatResponse{
		Chat: model.ChatToProto(chat),
	}, nil
}

func (h *Handler) CreateDirectChat(
	ctx context.Context,
	req *chatspb.CreateDirectChatRequest,
) (*chatspb.ChatResponse, error) {
	if req == nil {
		return nil, ErrRequestIsEmpty
	}

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, chats.ErrUserID.Error())
	}

	chat := &model.Chat{Type: model.ChatTypeDirect}
	log.Printf("Chat: %v\t", chat)
	chat, av, err := h.ctrl.CreateDirectChat(ctx, chat, uuid.MustParse(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	avProto, err := model.AvatarToProto(av)
	if err != nil {
		log.Printf("Handler:CreateDirectChat: %v", err)
		return &chatspb.ChatResponse{
			Chat: model.ChatToProto(chat),
		}, nil

	} else {
		return &chatspb.ChatResponse{
			Chat:   model.ChatToProto(chat),
			Avatar: avProto,
		}, nil
	}
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

	log.Printf("Update chat: %v\n", req.Chat)
	chat, err := model.ChatFromProto(req.Chat)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var avatar_data []byte
	if len(req.Avatar.Data) != 0 {
		avatar_data = req.Avatar.Data
	}

	err = h.ctrl.UpdateChat(ctx, chat, avatar_data)
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
		pProto, err := model.ChatParticipantToProto(&p)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		resp.Users = append(resp.Users, pProto)
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

func (h *Handler) ListChatsForUser(ctx context.Context, req *chatspb.ListChatsForUserRequest) (
	*chatspb.ListChatsForUserResponse,
	error,
) {
	if req == nil {
		return nil, ErrRequestIsEmpty
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	chats, err := h.ctrl.ListChatsForUser(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &chatspb.ListChatsForUserResponse{}
	for _, c := range chats {
		resp.Chats = append(resp.Chats, model.ChatToProto(&c))
	}

	return resp, nil
}

func (h *Handler) ListAvatarsForChats(ctx context.Context, req *chatspb.ListAvatarsForChatsRequest) (
	*chatspb.ListAvatarsForChatsResponse,
	error,
) {
	if req == nil {
		return nil, ErrRequestIsEmpty
	}

	var group_ids []uuid.UUID
	var direct_ids []uuid.UUID
	for _, chatSummary := range req.Chats {
		id, err := uuid.Parse(chatSummary.Id)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		if chatSummary.IsDirect {
			direct_ids = append(direct_ids, id)
		} else {
			group_ids = append(group_ids, id)
		}
	}

	group_chats_avs, err := h.ctrl.ListAvatarsForGroupChats(ctx, group_ids)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	log.Printf("group chats len: %v\n", len(group_chats_avs))

	direct_chats_avs, err := h.ctrl.ListAvatarsForDirectChats(ctx, direct_ids)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	log.Printf("direct chats len: %v\n", len(direct_chats_avs))

	avatars := append(group_chats_avs, direct_chats_avs...)
	log.Printf("all chats len: %v\n", len(avatars))
	resp := &chatspb.ListAvatarsForChatsResponse{}
	for _, a := range avatars {
		aProto, err := model.AvatarToProto(&a)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		resp.Avatars = append(resp.Avatars, aProto)
	}

	return resp, nil
}
