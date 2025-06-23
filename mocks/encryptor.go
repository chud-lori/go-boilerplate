package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockEncryptor is a mock implementation of the Encryptor interface
// using testify/mock.
type MockEncryptor struct {
	mock.Mock
}

// HashPassword mocks the HashPassword method of the Encryptor interface.
// It records the call and returns the values configured by the expectations.
func (m *MockEncryptor) HashPassword(password string) (string, error) {
	args := m.Called(password)
	var r0 string
	if args.Get(0) != nil {
		r0 = args.Get(0).(string)
	}
	return r0, args.Error(1)
}

// CompareHash mocks the CompareHash method of the Encryptor interface.
// It records the call and returns the error configured by the expectations.
func (m *MockEncryptor) CompareHash(hash, password string) error {
	args := m.Called(hash, password)
	return args.Error(0)
}
