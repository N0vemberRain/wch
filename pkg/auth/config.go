package auth

import (
	"errors"
	"os"
	"strconv"
	"time"
)

var ErrTokenSecretUndefined = errors.New("JWT secret is undefined")
var ErrTokenExpireTimeUndefined = errors.New("JWT expire time is undefined")
var ErrTokenExpireTimeFormat = errors.New("JWT expire time must be integer")

type Config struct {
	Secret     string
	ExpireTime time.Duration
}

func LoadConfig() (*Config, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, ErrTokenSecretUndefined
	}

	expireTimeStr := os.Getenv("JWT_DURATION")
	if expireTimeStr == "" {
		return nil, ErrTokenExpireTimeUndefined
	}

	expireTime, err := strconv.Atoi(expireTimeStr)
	if err != nil {
		return nil, ErrTokenExpireTimeFormat
	}

	return &Config{
		Secret:     secret,
		ExpireTime: time.Duration(expireTime) * time.Minute,
	}, nil
}
