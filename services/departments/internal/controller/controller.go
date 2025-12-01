package controller

import (
	"context"
	"errors"

	mdmodel "wch/services/users/pkg/model"
)

var ErrNotFound = errors.New("userdata not found")

type userdataGateway interface {
	Get(ctx context.Context, id string) (*mdmodel.User, error)
}
