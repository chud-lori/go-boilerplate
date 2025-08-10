package queue

import (
	"context"
	"time"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

// RabbitMQJobQueue implements the JobQueue interface using RabbitMQ.
type RabbitMQJobQueue struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *logrus.Logger
}

func NewRabbitMQJobQueue(amqpURL string, logger *logrus.Logger) (ports.JobQueue, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &RabbitMQJobQueue{conn: conn, channel: ch, logger: logger}, nil
}

// PublishJob publishes a job to the specified queue (jobType).
func (q *RabbitMQJobQueue) PublishJob(ctx context.Context, jobType string, payload []byte) error {
	// Ensure the queue exists
	_, err := q.channel.QueueDeclare(
		jobType,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	return q.channel.PublishWithContext(
		ctx,
		"", // exchange
		jobType,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
			Timestamp:   time.Now(),
		},
	)
}

// ConsumeJobs consumes jobs from the specified queue (jobType) and calls handler for each message.
func (q *RabbitMQJobQueue) ConsumeJobs(ctx context.Context, jobType string, handler func([]byte) error) error {
	// Ensure the queue exists
	_, err := q.channel.QueueDeclare(
		jobType,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	msgs, err := q.channel.Consume(
		jobType,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}
			if err := handler(msg.Body); err != nil {
				q.logger.Errorf("Failed to handle job: %v", err)
			}
		}
	}
}

func (q *RabbitMQJobQueue) Close() error {
	q.channel.Close()
	return q.conn.Close()
}
