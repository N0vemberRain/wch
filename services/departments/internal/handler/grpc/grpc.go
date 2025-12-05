package grpc

import (
	"context"
	"wch/gen"

	"wch/services/departments/internal/controller"
	"wch/services/departments/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a controller gRPC handler.
type Handler struct {
	gen.UnimplementedDepsServiceServer
	ctrl *controller.Controller
}

// New creates a new user gRPC handler.
func New(dc *controller.Controller) *Handler {
	return &Handler{ctrl: dc}
}

func (h *Handler) GetDepByID(ctx context.Context, req *gen.GetDepByIDRequest) (*gen.GetDepByIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is empty")
	}

	dep, err := h.ctrl.GetByID(ctx, int(req.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.GetDepByIDResponse{
		Dep: model.DepToProto(dep),
	}, nil
}

func (h *Handler) GetDepByName(ctx context.Context, req *gen.GetDepByNameRequest) (*gen.GetDepByNameResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "req is empty")
	}

	dep, err := h.ctrl.GetByName(ctx, req.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.GetDepByNameResponse{Dep: model.DepToProto(dep)}, nil
}

func (h *Handler) ListDeps(ctx context.Context, req *gen.ListDepsRequest) (*gen.ListDepsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is empty")
	}

	deps, err := h.ctrl.GetAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &gen.ListDepsResponse{}
	for _, d := range deps {
		resp.Deps = append(resp.Deps, model.DepToProto(d))
	}

	return resp, nil
}
