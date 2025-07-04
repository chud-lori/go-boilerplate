package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MailClient struct {
	mock.Mock
}

func (_m *MailClient) SendMail(email, text string) error {
	args := _m.Called(email, text)
	return args.Error(0)
}
