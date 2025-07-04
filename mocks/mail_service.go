package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MailService struct {
	mock.Mock
}

func (_m *MailService) SendSignInNotification(ctx context.Context, email, text string) error {
	args := _m.Called(ctx, email, text)
	return args.Error(0)
}
