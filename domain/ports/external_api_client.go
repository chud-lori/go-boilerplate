package ports

import "context"

// ExternalApiClient defines the interface for making external API calls with circuit breaker
// This interface should be implemented by infrastructure/api_clients.ApiClient
// and injected into services that need to call external REST APIs.
type ExternalApiClient interface {
	DoRequest(ctx context.Context, method, url string, headers map[string]string, body []byte) ([]byte, error)
} 