package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestAPIKeyMiddleware_ValidKey(t *testing.T) {
	logger := logrus.New()
	called := false
	h := APIKeyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}), "secret", logger)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-KEY", "secret")
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if !called {
		t.Error("expected next handler to be called")
	}
}

func TestAPIKeyMiddleware_InvalidKey(t *testing.T) {
	logger := logrus.New()
	h := APIKeyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not call next handler")
	}), "secret", logger)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-KEY", "wrong")
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rw.Code)
	}
}

func TestAPIKeyMiddleware_MissingKey(t *testing.T) {
	logger := logrus.New()
	h := APIKeyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not call next handler")
	}), "secret", logger)

	req := httptest.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rw.Code)
	}
}

func TestAPIKeyMiddleware_SkipDocs(t *testing.T) {
	logger := logrus.New()
	called := false
	h := APIKeyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}), "secret", logger)

	req := httptest.NewRequest("GET", "/docs/index.html", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if !called {
		t.Error("expected next handler to be called for /docs/ path")
	}
} 