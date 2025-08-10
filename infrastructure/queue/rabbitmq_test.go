package queue_test

import (
    "context"
    "encoding/json"
    "testing"
    "time"

    "github.com/chud-lori/go-boilerplate/domain/entities"
    "github.com/chud-lori/go-boilerplate/infrastructure/queue"
    "github.com/chud-lori/go-boilerplate/internal/testutils"
    "github.com/google/uuid"
    "github.com/sirupsen/logrus"
    "github.com/stretchr/testify/assert"
)

func TestUploadJobMessage_MarshalUnmarshal(t *testing.T) {
	orig := entities.UploadJobMessage{
		UploadID:  "id1",
		PostID:    "pid",
		FileName:  "file.txt",
		FileType:  "text/plain",
		FileData:  []byte("data"),
		RequestID: "req-123",
	}
	data, err := json.Marshal(orig)
	assert.NoError(t, err)

	var out entities.UploadJobMessage
	assert.NoError(t, json.Unmarshal(data, &out))
	assert.Equal(t, orig.UploadID, out.UploadID)
	assert.Equal(t, orig.PostID, out.PostID)
	assert.Equal(t, orig.FileName, out.FileName)
	assert.Equal(t, orig.FileType, out.FileType)
	assert.Equal(t, orig.FileData, out.FileData)
	assert.Equal(t, orig.RequestID, out.RequestID)
}

func TestConsumeJobs_ValidAndInvalidPayload(t *testing.T) {
	// Simulate a valid and invalid payload
	validMsg := entities.UploadJobMessage{
		UploadID:  "id1",
		PostID:    "pid",
		FileName:  "file.txt",
		FileType:  "text/plain",
		FileData:  []byte("data"),
		RequestID: "req-123",
	}
	validData, _ := json.Marshal(validMsg)
	invalidData := []byte("not a json")

	ch := make(chan []byte, 2)
	ch <- validData
	ch <- invalidData
	close(ch)

	var handled []string
	handler := func(payload []byte) error {
		var msg entities.UploadJobMessage
		if err := json.Unmarshal(payload, &msg); err != nil {
			handled = append(handled, "invalid")
			return err
		}
		handled = append(handled, msg.UploadID)
		return nil
	}

	// Simulate ConsumeJobs logic
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				break
			}
			handler(msg)
		}
		if len(handled) == 2 {
			break
		}
	}

	assert.Contains(t, handled, "id1")     // Valid payload handled
	assert.Contains(t, handled, "invalid") // Invalid payload handled
}

func TestRabbitMQ_PublishConsume_Integration(t *testing.T) {
    t.Parallel()
	ctx := context.Background()

    // Use shared RabbitMQ container
    amqpURL, err := testutils.GetRabbitURL()
    assert.NoError(t, err)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

    // Give a brief moment for the consumer to attach before publishing
    time.Sleep(250 * time.Millisecond)

	jobQueue, err := queue.NewRabbitMQJobQueue(amqpURL, logger)
	assert.NoError(t, err)
	defer jobQueue.Close()

	queueName := "test_publish_consume"
	received := make(chan entities.UploadJobMessage, 1)

	// Start consumer
	go func() {
		jobQueue.ConsumeJobs(ctx, queueName, func(payload []byte) error {
			var msg entities.UploadJobMessage
			err := json.Unmarshal(payload, &msg)
			assert.NoError(t, err)
			received <- msg
			return nil
		})
	}()

	// Publish a message
	msg := entities.UploadJobMessage{
		UploadID:  uuid.New().String(),
		PostID:    uuid.New().String(),
		FileName:  "file.txt",
		FileType:  "text/plain",
		FileData:  []byte("data"),
		RequestID: "req-123",
	}
	payload, err := json.Marshal(msg)
	assert.NoError(t, err)

	err = jobQueue.PublishJob(ctx, queueName, payload)
	assert.NoError(t, err)

	// Wait for the message to be received
	select {
	case got := <-received:
		assert.Equal(t, msg.UploadID, got.UploadID)
		assert.Equal(t, msg.RequestID, got.RequestID)
		assert.Equal(t, msg.FileName, got.FileName)
	case <-time.After(5 * time.Second):
		t.Fatal("Did not receive message from RabbitMQ")
	}
}
