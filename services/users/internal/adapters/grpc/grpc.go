package grpc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	userspb "wch/gen/users/v1"
	//"wch/services/users/internal/controller"
	"wch/services/users/internal/controller"
	domain "wch/services/users/internal/domain"
	"wch/services/users/internal/domain/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a controller gRPC handler.
type Handler struct {
	userspb.UnimplementedUsersServiceServer
	usr_ctrl *controller.UserController
}

// New creates a new user gRPC handler.
func New(uc *controller.UserController) *Handler {
	return &Handler{usr_ctrl: uc}
}

func (h *Handler) CreateUser(ctx context.Context, req *userspb.CreateUserRequest) (*userspb.UserResponse, error) {
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

	return &userspb.UserResponse{}, nil
}

// GetUserByID returns user details by id.
func (h *Handler) GetUserByID(ctx context.Context, req *userspb.GetUserByIDRequest) (*userspb.GetUserResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	u, err := h.usr_ctrl.GetByID(ctx, req.UserId)
	if err != nil && errors.Is(err, domain.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &userspb.GetUserResponse{
		User: model.UserToProto(u),
	}, nil
}

func (h *Handler) GetUsersByIDs(ctx context.Context, req *userspb.GetUsersByIDsRequest) (
	*userspb.ListUsersResponse, error,
) {
	if len(req.Id) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	fmt.Println(len(req.Id))
	for _, id := range req.Id {
		fmt.Println(id)
	}

	users, err := h.usr_ctrl.GetByIDs(ctx, req.Id)
	if err != nil && errors.Is(err, domain.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &userspb.ListUsersResponse{}
	for _, u := range users {
		resp.Users = append(resp.Users, model.UserToProto(&u))
	}

	return resp, nil
}

// GetUserByEmail returns user details by id.
func (h *Handler) GetUserByEmail(ctx context.Context, req *userspb.GetUserByEmailRequest) (*userspb.GetUserResponse, error) {
	if req == nil || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty email")
	}
	u, err := h.usr_ctrl.GetByEmail(ctx, req.Email)
	if err != nil && errors.Is(err, domain.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &userspb.GetUserResponse{
		User: model.UserToProto(u),
	}, nil
}

func (h *Handler) SearchUsers(ctx context.Context, req *userspb.SearchUsersRequest) (*userspb.SearchUsersResponse, error) {
	filter := model.SearchFilter{
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

	resp := &userspb.SearchUsersResponse{}
	for _, u := range usersList {
		resp.Users = append(resp.Users, model.UserToProto(u))
	}

	return resp, nil
}

func (h *Handler) ListUsers(ctx context.Context, req *userspb.ListUsersRequest) (*userspb.ListUsersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is empty")
	}

	users, err := h.usr_ctrl.GetAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &userspb.ListUsersResponse{}
	for _, u := range users {
		resp.Users = append(resp.Users, model.UserToProto(u))
	}

	return resp, nil
}
