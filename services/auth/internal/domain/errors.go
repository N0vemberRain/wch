package domain

import "errors"

var (
	ErrCredentialsNotFound = errors.New("credentials not found")
	ErrInvalidPassword     = errors.New("password is invalid")
)
