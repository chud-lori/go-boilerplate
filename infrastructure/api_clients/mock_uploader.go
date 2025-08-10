package api_clients

import (
	"context"
	"fmt"
)

type MockUploader struct{}

func NewMockUploader() *MockUploader {
	return &MockUploader{}
}

// Upload simulates uploading a file to S3/MinIO by sleeping, then returns a mock URL.
func (u *MockUploader) Upload(ctx context.Context, fileName string, fileData []byte) (string, error) {
	// time.Sleep( * time.Second) // Simulate upload delay
	return fmt.Sprintf("https://mock-bucket.example.com/%s", fileName), nil
}
