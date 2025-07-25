package api_clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker/v2"
)

type ApiMailClient struct {
	Endpoint string
	Client   *http.Client
	Breaker  *gobreaker.CircuitBreaker[[]byte]
}

func defaultBreaker() *gobreaker.CircuitBreaker[[]byte] {
	var st gobreaker.Settings
	st.Name = "ApiMailClient"
	st.MaxRequests = 3
	st.Interval = 60 * time.Second
	st.Timeout = 10 * time.Second
	return gobreaker.NewCircuitBreaker[[]byte](st)
}

var _ ports.MailClient = (*ApiMailClient)(nil)

func NewApiMailClient(endpoint string) *ApiMailClient {
	return &ApiMailClient{
		Endpoint: endpoint,
		Client:   &http.Client{Timeout: 5 * time.Second},
		Breaker:  defaultBreaker(),
	}
}

type mailRequest struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

func (a *ApiMailClient) SendMail(ctx context.Context, email string, message string) error {
	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	_, err := a.Breaker.Execute(func() ([]byte, error) {
		payload := mailRequest{Email: email, Message: message}
		body, _ := json.Marshal(payload)
		req, err := http.NewRequestWithContext(ctx, "POST", a.Endpoint, bytes.NewBuffer(body))
		if err != nil {
			logger.WithError(err).Error("failed to create mail API request")
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := a.Client.Do(req)
		if err != nil {
			logger.WithError(err).Error("mail API request failed")
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			return nil, fmt.Errorf("mail API returned status: %d", resp.StatusCode)
		}
		logger.Infof("API mail sent to %s", email)
		return nil, nil
	})

	return err
}
