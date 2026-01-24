package ports

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) bool
}
