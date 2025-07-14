package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorsMiddleware_SetsHeaders(t *testing.T) {
	h := CorsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	headers := rw.Header()
	if headers.Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Error("missing or wrong Access-Control-Allow-Origin")
	}
	if headers.Get("Access-Control-Allow-Methods") == "" {
		t.Error("missing Access-Control-Allow-Methods")
	}
	if headers.Get("Access-Control-Allow-Headers") == "" {
		t.Error("missing Access-Control-Allow-Headers")
	}
	if headers.Get("Access-Control-Allow-Credentials") != "true" {
		t.Error("missing or wrong Access-Control-Allow-Credentials")
	}
}

func TestCorsMiddleware_HandlesOptions(t *testing.T) {
	h := CorsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not call next handler on OPTIONS")
	}))

	req := httptest.NewRequest("OPTIONS", "/", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("expected 200 for OPTIONS, got %d", rw.Code)
	}
} 