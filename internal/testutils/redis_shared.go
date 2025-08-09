package testutils

import (
	"context"
	"sync"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	redisOnce      sync.Once
	redisContainer testcontainers.Container
	redisAddr      string
	redisErr       error
	redisStopOnce  sync.Once
)

// StartRedisOnce ensures a single Redis container is started for the test process.
func StartRedisOnce() error {
	redisOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		req := testcontainers.ContainerRequest{
			Image:        "redis:7-alpine",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForListeningPort("6379/tcp").WithStartupTimeout(1 * time.Minute),
		}
		c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			redisErr = err
			return
		}
		redisContainer = c
		addr, err := c.Endpoint(ctx, "")
		if err != nil {
			redisErr = err
			return
		}
		redisAddr = addr
	})
	return redisErr
}

func GetRedisAddr() (string, error) {
	if err := StartRedisOnce(); err != nil {
		return "", err
	}
	return redisAddr, nil
}

func StopRedis() {
	redisStopOnce.Do(func() {
		if redisContainer != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			_ = redisContainer.Terminate(ctx)
		}
	})
}
