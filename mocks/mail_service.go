package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockMailService struct {
	mock.Mock
}

func (_m *MockMailService) SendSignInNotification(ctx context.Context, email, text string) error {
	args := _m.Called(ctx, email, text)
	return args.Error(0)
}
