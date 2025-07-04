package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockMailClient struct {
	mock.Mock
}

func (_m *MockMailClient) SendMail(ctx context.Context, email, text string) error {
	args := _m.Called(ctx, email, text)
	return args.Error(0)
}
