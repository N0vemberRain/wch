package users

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailExists     = errors.New("email already exists")
	ErrInvalidUserData = errors.New("invalid user's data")
)
