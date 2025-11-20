package controller

import (
	"context"
	"errors"

	users "wch/services/users/internal"
	"wch/services/users/internal/repository"
	"wch/services/users/pkg/model"

	"github.com/google/uuid"
)

// Controller defines an users service controller
type UserController struct {
	repo repository.UserRepository
}

func NewUserController(repo repository.UserRepository) *UserController {
	return &UserController{repo: repo}
}

func (c *UserController) CreateUser(ctx context.Context, u *model.User) error {
	user, _ := c.repo.GetUserByEmail(ctx, u.Email)
	if user != nil {
		return users.ErrEmailExists
	}

	u.ID = uuid.New()
	if err := c.repo.CreateUser(ctx, u); err != nil {
		return err
	}

	return nil
}

// Get returns the user's details
func (c *UserController) Get(ctx context.Context, id string) (*model.User, error) {
	userdata, err := c.repo.GetUserByID(ctx, id)
	if err != nil && errors.Is(err, errors.New("user not found")) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return userdata, nil
}
