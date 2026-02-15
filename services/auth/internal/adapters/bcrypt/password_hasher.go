package bcrypt

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {
	cost int
}

func NewPasswordHasher(cost int) *PasswordHasher {
	return &PasswordHasher{
		cost: cost,
	}
}

func (h *PasswordHasher) Hash(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), h.cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (h *PasswordHasher) Compare(hash, plain string) bool {
	fmt.Printf("hash %s, plain %s\n", hash, plain)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err != nil {
		return false
	}

	return true
}
