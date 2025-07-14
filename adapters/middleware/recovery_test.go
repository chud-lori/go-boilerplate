package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"bytes"

	"github.com/sirupsen/logrus"
)

func TestRecoveryMiddleware_RecoversFromPanic(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, DisableQuote: true})

	h := RecoveryMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(errors.New("test panic"))
	}), logger)

	req := httptest.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rw.Code)
	}
	if !strings.Contains(rw.Body.String(), "An unexpected error occurred") {
		t.Errorf("expected error message in response, got %q", rw.Body.String())
	}
	if !strings.Contains(buf.String(), "PANIC") {
		t.Errorf("expected panic to be logged, got %q", buf.String())
	}
} 