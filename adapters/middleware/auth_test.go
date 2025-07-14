package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chud-lori/go-boilerplate/mocks"
	"github.com/sirupsen/logrus"
)

func TestJWTMiddleware_MissingToken(t *testing.T) {
	logger := logrus.New()
	m := &mocks.MockTokenManager{}
	h := JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not call next handler")
	}), m, logger)

	req := httptest.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rw.Code)
	}
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	logger := logrus.New()
	m := &mocks.MockTokenManager{}
	m.On("ValidateToken", "badtoken").Return("", errors.New("invalid"))

	h := JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not call next handler")
	}), m, logger)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer badtoken")
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rw.Code)
	}
}

func TestJWTMiddleware_ValidToken(t *testing.T) {
	logger := logrus.New()
	m := &mocks.MockTokenManager{}
	m.On("ValidateToken", "goodtoken").Return("user123", nil)

	called := false
	h := JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		uid := r.Context().Value(UserIDKey)
		if uid != "user123" {
			t.Errorf("expected userID to be injected, got %v", uid)
		}
	}), m, logger)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer goodtoken")
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if !called {
		t.Error("expected next handler to be called")
	}
} 