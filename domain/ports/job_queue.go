package ports

import (
	"context"
)

type JobQueue interface {
	// PublishJob publishes a job to the queue. The payload should be a serialized message (e.g., JSON).
	PublishJob(ctx context.Context, jobType string, payload []byte) error
	// ConsumeJobs starts consuming jobs of a given type, calling handler for each message.
	ConsumeJobs(ctx context.Context, jobType string, handler func([]byte) error) error
	// Close closes any resources held by the queue.
	Close() error
} 