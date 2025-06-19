package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/sirupsen/logrus"
)

type contextKey string

const UserIDKey contextKey = "userID"

func JWTMiddleware(tokenManager ports.TokenManager, logger *logrus.Logger) func(http.Handler) http.Handler {
	mwLogger := logger.WithFields(logrus.Fields{
		"layer": "middleware",
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := tokenManager.ValidateToken(token)
			if err != nil {
				mwLogger.Warnf("Invalid token: %v", err)
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
