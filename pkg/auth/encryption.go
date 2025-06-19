package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptEncryptor struct{}

func (e *BcryptEncryptor) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (e *BcryptEncryptor) CompareHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
