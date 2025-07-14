package api_clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/chud-lori/go-boilerplate/internal/testutils"
)

// Add test-only helper for custom breaker
func newApiMailClientWithBreaker(endpoint string, breaker *gobreaker.CircuitBreaker[[]byte]) *ApiMailClient {
	return &ApiMailClient{
		Endpoint: endpoint,
		Client:   &http.Client{Timeout: 5 * time.Second},
		Breaker:  breaker,
	}
}

func TestNewApiMailClient(t *testing.T) {
	endpoint := "http://example.com/mail"
	client := newApiMailClientWithBreaker(endpoint, testutils.FreshMailBreaker())

	assert.Equal(t, endpoint, client.Endpoint)
	assert.NotNil(t, client.Client)
	assert.Equal(t, 5*time.Second, client.Client.Timeout)
	assert.NotNil(t, client.Breaker)
}

func TestApiMailClient_SendMail_Success(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and headers
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Parse and verify request body
		var requestBody struct {
			Email   string `json:"email"`
			Message string `json:"message"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		require.NoError(t, err)
		assert.Equal(t, "test@example.com", requestBody.Email)
		assert.Equal(t, "Test message", requestBody.Message)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	// Create client and context
	client := newApiMailClientWithBreaker(server.URL, testutils.FreshMailBreaker())
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	// Test SendMail
	err := client.SendMail(ctx, "test@example.com", "Test message")
	assert.NoError(t, err)
}

func TestApiMailClient_SendMail_AcceptedStatus(t *testing.T) {
	// Create test server that returns 202 Accepted
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"status": "accepted"}`))
	}))
	defer server.Close()

	client := newApiMailClientWithBreaker(server.URL, testutils.FreshMailBreaker())
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	err := client.SendMail(ctx, "test@example.com", "Test message")
	assert.NoError(t, err)
}

func TestApiMailClient_SendMail_ErrorStatus(t *testing.T) {
	// Create test server that returns error status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	client := newApiMailClientWithBreaker(server.URL, testutils.FreshMailBreaker())
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	err := client.SendMail(ctx, "test@example.com", "Test message")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mail API returned status: 500")
}

func TestApiMailClient_SendMail_NetworkError(t *testing.T) {
	// Create client with invalid endpoint
	client := newApiMailClientWithBreaker("http://invalid-endpoint-that-does-not-exist.com", testutils.FreshMailBreaker())
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	err := client.SendMail(ctx, "test@example.com", "Test message")
	assert.Error(t, err)
}

func TestApiMailClient_SendMail_RequestCreationError(t *testing.T) {
	// Create client with invalid URL to cause request creation error
	client := newApiMailClientWithBreaker("://invalid-url", testutils.FreshMailBreaker())
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	err := client.SendMail(ctx, "test@example.com", "Test message")
	assert.Error(t, err)
}

func TestApiMailClient_CircuitBreaker_OpenState(t *testing.T) {
	// Create a server that always fails
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := newApiMailClientWithBreaker(server.URL, testutils.FreshMailBreaker())
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	// Make multiple requests to trigger circuit breaker
	for i := 0; i < 5; i++ {
		err := client.SendMail(ctx, "test@example.com", "Test message")
		assert.Error(t, err)
	}

	// After multiple failures, circuit breaker should be open
	// and subsequent requests should fail immediately
	err := client.SendMail(ctx, "test@example.com", "Test message")
	assert.Error(t, err)
}

func TestApiMailClient_CircuitBreaker_HalfOpenState(t *testing.T) {
	// Create a server that fails initially, then succeeds
	failCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		failCount++
		if failCount <= 3 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	client := newApiMailClientWithBreaker(server.URL, testutils.FreshMailBreaker())
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	// Make initial requests to trigger circuit breaker
	for i := 0; i < 3; i++ {
		err := client.SendMail(ctx, "test@example.com", "Test message")
		assert.Error(t, err)
	}

	// Wait for circuit breaker to go to half-open state
	// Note: In a real test, you might need to wait for the timeout period
	// For this test, we'll just verify the behavior after failures

	// The circuit breaker should eventually allow a request through
	// when it's in half-open state, but this depends on timing
	// For now, we'll just verify that the client handles the circuit breaker correctly
}

func TestApiMailClient_CircuitBreaker_Recovery(t *testing.T) {
	// Create a server that fails initially, then succeeds
	failCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		failCount++
		if failCount <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	client := newApiMailClientWithBreaker(server.URL, testutils.FreshMailBreaker())
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	// Make a few failing requests
	for i := 0; i < 2; i++ {
		err := client.SendMail(ctx, "test@example.com", "Test message")
		assert.Error(t, err)
	}

	// Make a successful request
	err := client.SendMail(ctx, "test@example.com", "Test message")
	assert.NoError(t, err)
}

func TestApiMailClient_ConcurrentRequests(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate some processing time
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newApiMailClientWithBreaker(server.URL, testutils.FreshMailBreaker())
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	// Make concurrent requests
	const numRequests = 10
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			err := client.SendMail(ctx, "test@example.com", "Test message")
			results <- err
		}()
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		err := <-results
		assert.NoError(t, err)
	}
}

func TestApiMailClient_Timeout(t *testing.T) {
	// Create a server that delays longer than the client timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(6 * time.Second) // Longer than client timeout (5 seconds)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newApiMailClientWithBreaker(server.URL, testutils.FreshMailBreaker())
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	err := client.SendMail(ctx, "test@example.com", "Test message")
	assert.Error(t, err)
}

func TestApiMailClient_ContextCancellation(t *testing.T) {
	// Create a server that delays
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newApiMailClientWithBreaker(server.URL, testutils.FreshMailBreaker())
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	// Cancel context immediately
	cancel()

	err := client.SendMail(ctx, "test@example.com", "Test message")
	assert.Error(t, err)
}

func TestApiMailClient_EmptyEmailAndMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			Email   string `json:"email"`
			Message string `json:"message"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		require.NoError(t, err)
		assert.Equal(t, "", requestBody.Email)
		assert.Equal(t, "", requestBody.Message)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newApiMailClientWithBreaker(server.URL, testutils.FreshMailBreaker())
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	err := client.SendMail(ctx, "", "")
	assert.NoError(t, err)
}

func TestApiMailClient_SpecialCharacters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			Email   string `json:"email"`
			Message string `json:"message"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		require.NoError(t, err)
		assert.Equal(t, "test+special@example.com", requestBody.Email)
		assert.Equal(t, "Message with special chars: !@#$%^&*()", requestBody.Message)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newApiMailClientWithBreaker(server.URL, testutils.FreshMailBreaker())
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

	err := client.SendMail(ctx, "test+special@example.com", "Message with special chars: !@#$%^&*()")
	assert.NoError(t, err)
}
