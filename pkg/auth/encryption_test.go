package auth_test

import (
	"testing"

	"github.com/chud-lori/go-boilerplate/pkg/auth"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptEncryptor_HashAndCompare_Success(t *testing.T) {
	encryptor := &auth.BcryptEncryptor{}
	password := "MySecurePassword123!"

	hash, err := encryptor.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	err = encryptor.CompareHash(hash, password)
	assert.NoError(t, err) // Should succeed
}

func TestBcryptEncryptor_CompareHash_Failure(t *testing.T) {
	encryptor := &auth.BcryptEncryptor{}
	password := "correct-password"
	wrongPassword := "wrong-password"

	hash, err := encryptor.HashPassword(password)
	assert.NoError(t, err)

	err = encryptor.CompareHash(hash, wrongPassword)
	assert.ErrorIs(t, err, bcrypt.ErrMismatchedHashAndPassword)
}
