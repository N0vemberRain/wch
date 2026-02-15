package jwt

import (
	"errors"
	"fmt"
	"time"
	"wch/services/auth/internal/domain/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenIssuer struct {
	secret     []byte
	expireTime time.Duration
}

func NewTokenIssuer(secret string, expire time.Duration) *TokenIssuer {
	return &TokenIssuer{
		secret:     []byte(secret),
		expireTime: expire,
	}
}

func (iss *TokenIssuer) Issue(userID uuid.UUID) (*model.Token, error) {
	exp := time.Now().Add(iss.expireTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"role": "user",
		"exp":  exp.Unix(),
	})

	tokenStr, err := token.SignedString(iss.secret)
	if err != nil {
		return nil, err
	}

	return &model.Token{
		Value:     tokenStr,
		ExpiresAt: exp,
		UserID:    userID,
	}, nil
}

func (iss *TokenIssuer) Validate(tokenStr string) (*model.Token, error) {
	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("inexpected signing method: %v", token.Header["arg"])
		}

		return iss.secret, nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, errors.New("invalid token")
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

	return &model.Token{
		Value:     tokenStr,
		ExpiresAt: expTime.Time,
		UserID:    userID,
	}, nil
}
