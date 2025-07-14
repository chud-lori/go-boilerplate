package api_clients

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker/v2"
)

var _ ports.ExternalApiClient = (*ApiClient)(nil)

type ApiClient struct {
	Client  *http.Client
	Breaker *gobreaker.CircuitBreaker[[]byte]
	logger  *logrus.Entry
}

func NewApiClient(name string, logger *logrus.Logger) *ApiClient {
	st := gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     10 * time.Second,
	}

	apiCallLogger := logger.WithFields(logrus.Fields{
		"layer":  "ExternalApiCall",
		"driver": name,
	})

	return &ApiClient{
		Client:  &http.Client{Timeout: 5 * time.Second},
		Breaker: gobreaker.NewCircuitBreaker[[]byte](st),
		logger:  apiCallLogger,
	}
}

func (a *ApiClient) DoRequest(ctx context.Context, method, url string, headers map[string]string, body []byte) ([]byte, error) {
	return a.Breaker.Execute(func() ([]byte, error) {
		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := a.Client.Do(req)
		if err != nil {
			a.logger.WithError(err).Error("Failed do request")
			return nil, err
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("API returned status: %d, body: %s", resp.StatusCode, string(respBody))
		}
		return respBody, nil
	})
}
