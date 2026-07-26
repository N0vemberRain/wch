package controller

import (
	"context"
	"errors"
	"log"

	"wch/services/users/internal/domain"
	"wch/services/users/internal/domain/model"
	repository "wch/services/users/internal/domain/ports"

	"github.com/google/uuid"
)

// Controller defines an users service controller
type UserController struct {
	repo       repository.UserRepository
	av_storage repository.AvatarStorage
}

func NewUserController(
	repo repository.UserRepository,
	av_storage repository.AvatarStorage,
) *UserController {
	return &UserController{repo: repo, av_storage: av_storage}
}

func (c *UserController) CreateUser(ctx context.Context, u *model.User) error {
	user, _ := c.repo.GetUserByEmail(ctx, u.Email)
	if user != nil {
		return domain.ErrEmailExists
	}

	u.ID = uuid.New()
	if err := c.repo.CreateUser(ctx, u); err != nil {
		return err
	}

	return nil
}

func (c *UserController) UpdateUser(ctx context.Context, u *model.User, av_bytes []byte) error {
	log.Println("UserController.UpdateUser: enter")
	if u == nil {
		return domain.ErrInvalidUserData
	}
	log.Println("UserController.UpdateUser: before get user")
	oldUsr, err := c.repo.GetUserByID(ctx, u.ID.String())
	if err != nil {
		return err
	}

	log.Println("UserController.UpdateUser: before avatar saving")
	if len(av_bytes) != 0 {
		key, err := c.av_storage.Save(ctx, u.ID, av_bytes)
		if err != nil {
			return err
		}

		u.AvatarURL = key
	} else {
		log.Println("UserController.UpdateUser:  av bytes empty")
		u.AvatarURL = oldUsr.AvatarURL
	}

	log.Println("UserController.UpdateUser: end")
	return c.repo.UpdateUser(ctx, u)
}

// GetByID returns the user's details
func (c *UserController) GetByID(ctx context.Context, id string) (*model.User, error) {
	userdata, err := c.repo.GetUserByID(ctx, id)
	if err != nil && errors.Is(err, errors.New("user not found")) {
		return nil, domain.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return userdata, nil
}

func (c *UserController) GetByIDs(ctx context.Context, ids []string) ([]model.User, error) {
	return c.repo.GetUsersByIDs(ctx, ids)
}

// GetByEmail returns the user's details
func (c *UserController) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	userdata, err := c.repo.GetUserByEmail(ctx, email)
	if err != nil && errors.Is(err, errors.New("user not found")) {
		return nil, domain.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return userdata, nil
}

func (c *UserController) SearchUsers(ctx context.Context, filter model.SearchFilter) (
	[]*model.User, error,
) {
	return c.repo.SearchUsers(ctx, filter)
}

func (c *UserController) GetAll(ctx context.Context) ([]*model.User, error) {
	return c.repo.GetAll(ctx)
}

func (c *UserController) GetAvatarForUser(ctx context.Context, user_id string) ([]byte, error) {
	u, err := c.repo.GetUserByID(ctx, user_id)
	if err != nil {
		return nil, err
	}

	return c.av_storage.Get(ctx, u.AvatarURL)
}

func (c *UserController) GetAvatarsForChats(ctx context.Context, ids []uuid.UUID) (
	[]model.Avatar,
	error,
) {
	var avatars []model.Avatar
	for _, id := range ids {
		bytes, err := c.av_storage.Get(ctx, id.String())
		if err != nil {
			break
		}

		var a model.Avatar
		a.Data = bytes
		a.MimeType = "PNG"
		a.OwnerID = id

		avatars = append(avatars, a)
	}

	return avatars, nil
}
