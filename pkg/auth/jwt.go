package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	SecretKey  string
	Expiration time.Duration
}

func (j *JWTManager) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(j.Expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWTManager) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id not found in token")
	}

	return userID, nil
}
