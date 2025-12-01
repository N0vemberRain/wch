package grpc

import (
	"context"
	"errors"
	"strings"

	"wch/gen"
	"wch/pkg/models"

	//"wch/services/users/internal/controller"
	users "wch/services/users/internal"
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

// GetUserByID returns user details by id.
func (h *Handler) GetUserByID(ctx context.Context, req *gen.GetUserByIDRequest) (*gen.GetUserResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	u, err := h.usr_ctrl.GetByID(ctx, req.UserId)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.GetUserResponse{
		User: model.UserToProto(u),
	}, nil
}

// GetUserByEmail returns user details by id.
func (h *Handler) GetUserByEmail(ctx context.Context, req *gen.GetUserByEmailRequest) (*gen.GetUserResponse, error) {
	if req == nil || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty email")
	}
	u, err := h.usr_ctrl.GetByEmail(ctx, req.Email)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.GetUserResponse{
		User: model.UserToProto(u),
	}, nil
}

func (h *Handler) SearchUsers(ctx context.Context, req *gen.SearchUsersRequest) (*gen.SearchUsersResponse, error) {
	filter := users.SearchFilter{
		Email:        req.Email,
		Username:     req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Surname:      req.Surname,
		DepartmentID: int(req.DepartmentId),
	}

	usersList, err := h.usr_ctrl.SearchUsers(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	resp := &gen.SearchUsersResponse{}
	for _, u := range usersList {
		resp.Users = append(resp.Users, model.UserToProto(u))
	}

	return resp, nil
}

func (h *Handler) GetDepByID(ctx context.Context, req *gen.GetDepByIDRequest) (*gen.GetDepByIDResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	dep, err := h.dep_ctrl.GetByID(ctx, int(req.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.GetDepByIDResponse{
		Dep: models.DepToProto(dep),
	}, nil
}

func (h *Handler) GetDepByName(ctx context.Context, req *gen.GetDepByNameRequest) (*gen.GetDepByNameResponse, error) {
	if req == nil || req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil request or empty name")
	}

	dep, err := h.dep_ctrl.GetByName(ctx, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetDepByNameResponse{
		Dep: models.DepToProto(dep),
	}, nil
}

func (h *Handler) ListDeps(ctx context.Context, req *gen.ListDepsRequest) (*gen.ListDepsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "requset is empty")
	}

	deps, err := h.dep_ctrl.GetAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &gen.ListDepsResponse{}
	for _, d := range deps {
		resp.Deps = append(resp.Deps, models.DepToProto(d))
	}

	return resp, nil
}
