package api_clients

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"errors"
	"github.com/sony/gobreaker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/chud-lori/go-boilerplate/internal/testutils"
)

func TestApiClient_DoRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"ok"}`))
	}))
	defer server.Close()

	client := &ApiClient{
		Client:  &http.Client{Timeout: 2 * time.Second},
		Breaker: testutils.FreshBreaker("TestApiClientSuccess"),
	}

	ctx := context.Background()
	resp, err := client.DoRequest(ctx, "POST", server.URL, map[string]string{"Content-Type": "application/json"}, []byte(`{"foo":"bar"}`))
	assert.NoError(t, err)
	assert.Contains(t, string(resp), "ok")
}

func TestApiClient_DoRequest_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"fail"}`))
	}))
	defer server.Close()

	client := &ApiClient{
		Client:  &http.Client{Timeout: 2 * time.Second},
		Breaker: testutils.FreshBreaker("TestApiClientError"),
	}

	ctx := context.Background()
	resp, err := client.DoRequest(ctx, "GET", server.URL, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestApiClient_DoRequest_CircuitBreakerOpen(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	breaker := testutils.FreshBreaker("TestApiClientBreaker")
	client := &ApiClient{
		Client:  &http.Client{Timeout: 2 * time.Second},
		Breaker: breaker,
	}

	ctx := context.Background()
	// Trigger failures to open the breaker
	for i := 0; i < 10; i++ {
		_, _ = client.DoRequest(ctx, "GET", server.URL, nil, nil)
	}
	// Now the breaker should be open, and requests should fail fast
	_, err := client.DoRequest(ctx, "GET", server.URL, nil, nil)
	if err == nil {
		t.Fatalf("expected error when circuit breaker is open, got nil")
	}
	if !errors.Is(err, gobreaker.ErrOpenState) && !strings.Contains(err.Error(), "open") {
		t.Errorf("expected circuit breaker open error, got: %v", err)
	}
} 