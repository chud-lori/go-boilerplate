package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockExternalApiClient struct {
	mock.Mock
}

func (m *MockExternalApiClient) DoRequest(ctx context.Context, method, url string, headers map[string]string, body []byte) ([]byte, error) {
	args := m.Called(ctx, method, url, headers, body)
	var r0 []byte
	if args.Get(0) != nil {
		r0 = args.Get(0).([]byte)
	}
	return r0, args.Error(1)
}
