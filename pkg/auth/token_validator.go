package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Token struct {
	Value     string
	ExpiresAt time.Time
	UserID    uuid.UUID
}

type TokenValidator struct {
	secret     []byte
	expireTime time.Duration
}

func NewTokenValidator(secret string, expire time.Duration) *TokenValidator {
	return &TokenValidator{
		secret:     []byte(secret),
		expireTime: expire,
	}
}

func (tv *TokenValidator) Validate(tokenStr string) (*Token, error) {
	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("inexpected signing method: %v", token.Header["arg"])
		}

		return tv.secret, nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token: %s\n", err.Error())
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	expTime, err := claims.GetExpirationTime()
	if !ok {
		return nil, errors.New("expires time claim not found")
	}
	userStr, err := claims.GetSubject()
	if err != nil {
		return nil, errors.New("user id not found")
	}
	userID, err := uuid.Parse(userStr)
	if err != nil {
		return nil, errors.New("user id is invalid")
	}

	return &Token{
		Value:     tokenStr,
		ExpiresAt: expTime.Time,
		UserID:    userID,
	}, nil
}
