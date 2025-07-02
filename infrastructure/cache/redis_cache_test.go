package cache_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chud-lori/go-boilerplate/infrastructure/cache"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
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
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logrus.NewEntry(logrus.New()))

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

	// --- New Test for InvalidateByPrefix ---
	t.Run("InvalidateByPrefix", func(t *testing.T) {
		prefix := "posts:"
		keys := []string{
			prefix + "search=keyword1:page=1:limit=10",
			prefix + "search=keyword1:page=2:limit=10",
			prefix + "search=keyword2:page=1:limit=10",
			prefix + "all:page=1:limit=10", // Example of a non-search key with the prefix
			"other-key",                    // Should not be affected
		}
		values := []byte("some-post-data")
		exp := 5 * time.Minute

		// Set multiple keys with the prefix
		for _, key := range keys {
			err := cache.Set(ctx, key, values, exp)
			assert.NoError(t, err, fmt.Sprintf("Failed to set key: %s", key))
		}

		// Verify they are set
		for _, key := range keys {
			if key == "other-key" {
				continue // Skip the non-prefixed key for this check
			}
			got, err := cache.Get(ctx, key)
			assert.NoError(t, err, fmt.Sprintf("Failed to get key after set: %s", key))
			assert.Equal(t, string(values), got, fmt.Sprintf("Value mismatch for key: %s", key))
		}

		// Invalidate by prefix
		err := cache.InvalidateByPrefix(ctx, prefix)
		assert.NoError(t, err)

		// Verify keys with the prefix are deleted
		for _, key := range keys {
			if key == "other-key" {
				continue // Skip the non-prefixed key
			}
			got, err := cache.Get(ctx, key)
			assert.NoError(t, err, fmt.Sprintf("Error getting key after invalidation: %s", key))
			assert.Empty(t, got, fmt.Sprintf("Key '%s' should have been deleted but still found", key))
		}

		// Verify the non-prefixed key is still present
		gotOther, err := cache.Get(ctx, "other-key")
		assert.NoError(t, err)
		assert.Equal(t, string(values), gotOther, "'other-key' should not have been affected")
	})

	t.Run("Close", func(t *testing.T) {
		err := cache.Close()
		assert.NoError(t, err)
	})
}
