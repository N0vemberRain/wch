package controller

import (
	"context"
	"errors"
	"wch/services/auth/internal/domain"
	"wch/services/auth/internal/domain/model"
	"wch/services/auth/internal/domain/ports"
)

type Controller struct {
	repo   ports.CredentialsRepository
	hasher ports.PasswordHasher
	issuer ports.TokenIssuer
}

func NewAuthController(
	repo ports.CredentialsRepository,
	hasher ports.PasswordHasher,
	issuer ports.TokenIssuer,
) *Controller {
	return &Controller{
		repo:   repo,
		hasher: hasher,
		issuer: issuer,
	}
}

func (c *Controller) Login(ctx context.Context, email, password string) (*model.Token, error) {
	user, err := c.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, domain.ErrCredentialsNotFound
	}

	if password == "" || user.PasswordHash == "" {
		return nil, errors.New("password is empty")
	}

	if !c.hasher.Compare(user.PasswordHash, password) {
		return nil, domain.ErrInvalidPassword
	}

	token, err := c.issuer.Issue(user.UserID)
	if err != nil {
		return nil, err
	}

	return token, nil
}
