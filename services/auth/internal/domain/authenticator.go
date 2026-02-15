package domain

import "wch/services/auth/internal/domain/ports"

type Authenticator struct {
	credRepo ports.CredentialsRepository
	hasher   ports.PasswordHasher
	issuer   ports.TokenIssuer
}

func NewAuthenticator(
	repo ports.CredentialsRepository,
	hasher ports.PasswordHasher,
	issuer ports.TokenIssuer,
) *Authenticator {
	return &Authenticator{
		credRepo: repo,
		hasher:   hasher,
		issuer:   issuer,
	}
}
