package auth_test

import (
	"errors"
	"testing"
	"time"

	"github.com/chud-lori/go-boilerplate/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTManager_GenerateAndValidateToken(t *testing.T) {
	manager := &auth.JWTManager{
		SecretKey:  "supersecret",
		Expiration: time.Minute,
	}

	userID := "user123"

	tokenStr, err := manager.GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	parsedID, err := manager.ValidateToken(tokenStr)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedID)
}

func TestJWTManager_ValidateToken_InvalidSignature(t *testing.T) {
	// Token generated with a different secret
	otherManager := &auth.JWTManager{
		SecretKey:  "wrongsecret",
		Expiration: time.Minute,
	}
	manager := &auth.JWTManager{
		SecretKey:  "correctsecret",
		Expiration: time.Minute,
	}

	tokenStr, _ := otherManager.GenerateToken("user123")

	_, err := manager.ValidateToken(tokenStr)
	assert.Error(t, err)
}

func TestJWTManager_ValidateToken_ExpiredToken(t *testing.T) {
	manager := &auth.JWTManager{
		SecretKey:  "secret",
		Expiration: -time.Minute, // already expired
	}

	tokenStr, err := manager.GenerateToken("user123")
	assert.NoError(t, err)

	_, err = manager.ValidateToken(tokenStr)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, jwt.ErrTokenExpired) || err.Error() == "token is expired")
}

func TestJWTManager_ValidateToken_MissingUserID(t *testing.T) {
	// Manually create a token with no user_id
	manager := &auth.JWTManager{
		SecretKey:  "secret",
		Expiration: time.Minute,
	}

	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(manager.SecretKey))
	assert.NoError(t, err)

	_, err = manager.ValidateToken(tokenStr)
	assert.EqualError(t, err, "user_id not found in token")
}
