package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockMailClient struct {
	mock.Mock
}

func (_m *MockMailClient) SendMail(email, text string) error {
	args := _m.Called(email, text)
	return args.Error(0)
}
