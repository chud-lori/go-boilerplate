package mocks // Or a suitable test package, e.g., your_project_path/ports/mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockCache is a mock implementation of the ports.Cache interface.
type MockCache struct {
	mock.Mock
}

// Get mocks the Get method of the ports.Cache interface.
func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

// Set mocks the Set method of the ports.Cache interface.
func (m *MockCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

// Delete mocks the Delete method of the ports.Cache interface.
func (m *MockCache) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// Close mocks the Close method of the ports.Cache interface.
func (m *MockCache) Close() error {
	args := m.Called()
	return args.Error(0)
}

// ---

// For the NewRedisCache function, you typically wouldn't mock the constructor
// itself, but rather mock the dependencies it *uses* (like redis.Client if
// you were building a mock for NewRedisCache's *internal* behavior, which is
// less common).
//
// Instead, you'd inject the ports.Cache interface into your services,
// and in your tests, you'd inject an instance of MockCache.
