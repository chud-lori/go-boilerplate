package testutils

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

var (
    rabbitOnce      sync.Once
    rabbitContainer testcontainers.Container
    rabbitURL       string
    rabbitErr       error
    rabbitStopOnce  sync.Once
)

// StartRabbitOnce ensures a single RabbitMQ container is started for the test process.
func StartRabbitOnce() error {
    rabbitOnce.Do(func() {
        ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
        defer cancel()
        req := testcontainers.ContainerRequest{
            Image:        "rabbitmq:3.12-management-alpine",
            ExposedPorts: []string{"5672/tcp"},
            WaitingFor:   wait.ForLog("Server startup complete").WithStartupTimeout(2 * time.Minute),
        }
        c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
            ContainerRequest: req,
            Started:          true,
        })
        if err != nil {
            rabbitErr = err
            return
        }
        rabbitContainer = c

        host, err := c.Host(ctx)
        if err != nil {
            rabbitErr = err
            return
        }
        port, err := c.MappedPort(ctx, "5672")
        if err != nil {
            rabbitErr = err
            return
        }
        rabbitURL = fmt.Sprintf("amqp://guest:guest@%s:%s/", host, port.Port())
    })
    return rabbitErr
}

func GetRabbitURL() (string, error) {
    if err := StartRabbitOnce(); err != nil {
        return "", err
    }
    return rabbitURL, nil
}

func StopRabbit() {
    rabbitStopOnce.Do(func() {
        if rabbitContainer != nil {
            ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
            defer cancel()
            _ = rabbitContainer.Terminate(ctx)
        }
    })
}


