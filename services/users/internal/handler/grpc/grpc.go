package grpc

import (
	"context"
	"errors"
	"strings"

	"wch/gen"
	//"wch/services/users/internal/controller"
	"wch/services/users/internal/controller"
	"wch/services/users/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a controller gRPC handler.
type Handler struct {
	gen.UnimplementedUsersServiceServer
	usr_ctrl *controller.UserController
	dep_ctrl *controller.DepartmentController
}

// New creates a new user gRPC handler.
func New(uc *controller.UserController, dc *controller.DepartmentController) *Handler {
	return &Handler{usr_ctrl: uc, dep_ctrl: dc}
}

func (h *Handler) CreateUser(ctx context.Context, req *gen.CreateUserRequest) (*gen.UserResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	} else if req.User.Username == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username is empty")
	} else if req.User.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is empty")
	} else if !strings.Contains(req.User.Email, "@") {
		return nil, status.Errorf(codes.InvalidArgument, "email must contain @")
	} else if req.User.PasswordHash == "" {
		return nil, status.Errorf(codes.InvalidArgument, "password hash is empty")
	}

	usr, err := model.UserFromProto(req.User)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = h.usr_ctrl.CreateUser(ctx, usr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.UserResponse{}, nil
}

// GetUserDetails returns user details by id.
func (h *Handler) GetUser(ctx context.Context, req *gen.GetUserRequest) (*gen.GetUserResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	u, err := h.usr_ctrl.Get(ctx, req.UserId)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.GetUserResponse{
		User: model.UserToProto(u),
	}, nil
}
