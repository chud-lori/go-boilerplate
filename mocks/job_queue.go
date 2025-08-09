package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockJobQueue is a mock implementation of the ports.JobQueue interface.
type MockJobQueue struct {
	mock.Mock
}

func (m *MockJobQueue) PublishJob(ctx context.Context, jobType string, payload []byte) error {
	args := m.Called(ctx, jobType, payload)
	return args.Error(0)
}

func (m *MockJobQueue) ConsumeJobs(ctx context.Context, jobType string, handler func([]byte) error) error {
	args := m.Called(ctx, jobType, handler)
	return args.Error(0)
}

func (m *MockJobQueue) Close() error {
	args := m.Called()
	return args.Error(0)
} 