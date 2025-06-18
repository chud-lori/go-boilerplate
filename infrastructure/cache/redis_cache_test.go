package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/chud-lori/go-boilerplate/infrastructure/cache"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupRedisContainer(ctx context.Context) (testcontainers.Container, string, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp").WithStartupTimeout(10 * time.Second),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		redisC.Terminate(ctx)
		return nil, "", err
	}

	return redisC, endpoint, nil
}

func TestRedisCache(t *testing.T) {
	ctx := context.Background()

	// Start Redis container
	redisC, addr, err := setupRedisContainer(ctx)
	assert.NoError(t, err)
	defer redisC.Terminate(ctx)

	logger := logrus.New()

	cache, err := cache.NewRedisCache(addr, "", 0, logger)
	assert.NoError(t, err)
	assert.NotNil(t, cache)

	t.Run("Set and Get", func(t *testing.T) {
		key := "test-key"
		value := []byte("hello")
		exp := 10 * time.Second

		err := cache.Set(ctx, key, value, exp)
		assert.NoError(t, err)

		got, err := cache.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, "hello", got)
	})

	t.Run("Delete", func(t *testing.T) {
		key := "test-delete"
		value := []byte("bye")
		_ = cache.Set(ctx, key, value, 10*time.Second)

		err := cache.Delete(ctx, key)
		assert.NoError(t, err)

		got, err := cache.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, "", got)
	})

	t.Run("Close", func(t *testing.T) {
		err := cache.Close()
		assert.NoError(t, err)
	})
}
