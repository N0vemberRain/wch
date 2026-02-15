package grpc

import (
	"context"
	"errors"
	authpb "wch/gen/auth/v1"
	"wch/services/auth/internal/controller"
	"wch/services/auth/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	timeconv "google.golang.org/protobuf/types/known/timestamppb"
)

var ErrEmptyRequest = errors.New("request is empty")

type Handler struct {
	authpb.UnimplementedAuthServiceServer
	ctrl *controller.Controller
}

func NewHandler(c *controller.Controller) *Handler {
	return &Handler{
		ctrl: c,
	}
}

func (h *Handler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, ErrEmptyRequest.Error())
	}

	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is empty")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is empty")
	}

	token, err := h.ctrl.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch err {
		case domain.ErrCredentialsNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case domain.ErrInvalidPassword:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	res := &authpb.LoginResponse{
		AccessToken: &authpb.Token{},
	}
	res.AccessToken.Value = token.Value
	res.AccessToken.UserId = token.UserID.String()
	res.AccessToken.ExpiresAt = timeconv.New(token.ExpiresAt)

	return res, nil
}
