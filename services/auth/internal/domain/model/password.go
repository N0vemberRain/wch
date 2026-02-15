package model

type Password struct {
	hash string
}

func PasswordFromPlain(plain string) (*Password, error) {
	return nil, nil
}

func (p Password) Compare(plain string) bool {
	return false
}
