package repository

import (
	"context"

	"wch/services/users/internal/domain/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUsersByIDs(ctx context.Context, ids []string) ([]model.User, error)
	GetUserByName(ctx context.Context, username string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, u *model.User) error
	DeleteUser(ctx context.Context, id string) error
	SearchUsers(ctx context.Context, filter model.SearchFilter) ([]model.UserSummary, error)
	GetAll(ctx context.Context) ([]*model.User, error)
}
