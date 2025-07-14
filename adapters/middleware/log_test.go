package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogTrafficMiddleware_LogsRequestAndResponse(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, DisableQuote: true})

	h := LogTrafficMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("ok"))
	}), logger)

	req := httptest.NewRequest("GET", "/test", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusTeapot {
		t.Errorf("expected 418, got %d", rw.Code)
	}
	if !strings.Contains(buf.String(), "Processed request") {
		t.Errorf("expected log entry, got %q", buf.String())
	}
}

func TestLogTrafficMiddleware_MasksSensitiveFields(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, DisableQuote: true})

	body := map[string]interface{}{"password": "secret", "foo": "bar"}
	b, _ := json.Marshal(body)

	h := LogTrafficMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), logger)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	logOutput := buf.String()
	if strings.Contains(logOutput, "password") {
		t.Errorf("expected password to be masked, got %q", logOutput)
	}
	if !strings.Contains(logOutput, "foo") {
		t.Errorf("expected foo to be logged, got %q", logOutput)
	}
} 