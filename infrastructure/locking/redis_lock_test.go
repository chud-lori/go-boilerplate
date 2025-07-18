package locking_test

import (
	"context"
	"testing"
	"time"

	"github.com/chud-lori/go-boilerplate/infrastructure/locking"
	"github.com/chud-lori/go-boilerplate/internal/testutils"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRedisLocker(t *testing.T) {
	ctx := context.Background() // Use plain context.Background for tests unless specific context needed for logs

	// Start Redis container
	redisC, addr, err := testutils.SetupRedisContainer(ctx)
	assert.NoError(t, err)
	defer redisC.Terminate(ctx)

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Initialize the RedisLocker
	locker, err := locking.NewRedisLocker(addr, "", 0, logger)
	assert.NoError(t, err)
	assert.NotNil(t, locker)

	t.Run("Acquire and Release Lock Successfully", func(t *testing.T) {
		lockKey := "test-lock-1"
		ttl := 5 * time.Second

		// Client A acquires the lock
		acquiredA, uniqueValueA, err := locker.AcquireLock(lockKey, ttl)
		assert.NoError(t, err)
		assert.True(t, acquiredA, "Client A should acquire the lock")
		assert.NotEmpty(t, uniqueValueA, "Unique value should not be empty")

		// Verify the key exists in Redis (optional, for deeper confidence)
		redisClient := redis.NewClient(&redis.Options{Addr: addr})
		val, err := redisClient.Get(context.Background(), lockKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, uniqueValueA, val, "Redis key value should match the unique value used by Client A")

		// Client A releases the lock
		err = locker.ReleaseLock(lockKey, uniqueValueA)
		assert.NoError(t, err, "Client A should successfully release its lock")

		// Verify the key is deleted from Redis
		_, err = redisClient.Get(context.Background(), lockKey).Result()
		assert.Equal(t, redis.Nil, err, "Lock key should be deleted after release")
	})

	t.Run("Cannot Acquire Lock When Already Held", func(t *testing.T) {
		lockKey := "test-lock-2"
		ttl := 5 * time.Second

		// Client A acquires the lock
		acquiredA, uniqueValueA, err := locker.AcquireLock(lockKey, ttl)
		assert.NoError(t, err)
		assert.True(t, acquiredA, "Client A should acquire the lock")
		assert.NotEmpty(t, uniqueValueA, "Unique value should not be empty for Client A")

		// Client B tries to acquire the same lock
		acquiredB, uniqueValueB, err := locker.AcquireLock(lockKey, ttl)
		assert.NoError(t, err)
		assert.False(t, acquiredB, "Client B should NOT acquire the lock (already held)")
		assert.Empty(t, uniqueValueB, "Unique value should be empty for failed acquisition")

		// Ensure Client A can still release its lock
		err = locker.ReleaseLock(lockKey, uniqueValueA)
		assert.NoError(t, err, "Client A should successfully release its lock")

		// Verify Client B cannot release a lock it didn't acquire
		err = locker.ReleaseLock(lockKey, uniqueValueB) // uniqueValueB should be ""
		assert.NoError(t, err, "Releasing with an empty unique value should not return an error but be a no-op / warning")
	})

	t.Run("Lock Expires Automatically", func(t *testing.T) {
		lockKey := "test-lock-3"
		ttl := 1 * time.Second // Short TTL for testing expiration

		acquired, uniqueValue, err := locker.AcquireLock(lockKey, ttl)
		assert.NoError(t, err)
		assert.True(t, acquired, "Should acquire the lock initially")
		assert.NotEmpty(t, uniqueValue)

		time.Sleep(ttl + 500*time.Millisecond) // Wait for lock to expire, plus a small buffer

		// Verify the key is deleted from Redis after expiration
		redisClient := redis.NewClient(&redis.Options{Addr: addr})
		_, err = redisClient.Get(context.Background(), lockKey).Result()
		assert.Equal(t, redis.Nil, err, "Lock key should be automatically deleted after expiration")

		// Attempting to release an expired lock should fail gracefully
		err = locker.ReleaseLock(lockKey, uniqueValue)
		assert.NoError(t, err, "Releasing an expired lock should not return an error from the locker, but indicate it wasn't owned")
		// The warning log in your ReleaseLock implementation should indicate this
	})

	t.Run("Cannot Release Lock With Wrong Unique Value (Stale Lock Protection)", func(t *testing.T) {
		lockKey := "test-lock-4"
		ttl := 1 * time.Second // Short TTL to simulate expiration quickly

		// Client A acquires the lock
		acquiredA, uniqueValueA, err := locker.AcquireLock(lockKey, ttl)
		assert.NoError(t, err)
		assert.True(t, acquiredA, "Client A should acquire the lock")

		// Simulate Client A's operation taking too long, causing the lock to expire
		time.Sleep(ttl + 500*time.Millisecond)

		// Client B acquires the same lock (it's now expired from Client A)
		acquiredB, uniqueValueB, err := locker.AcquireLock(lockKey, ttl)
		assert.NoError(t, err)
		assert.True(t, acquiredB, "Client B should acquire the lock after Client A's expired")
		assert.NotEqual(t, uniqueValueA, uniqueValueB, "Unique values must be different")

		// Client A (which now holds a stale uniqueValueA) attempts to release the lock
		err = locker.ReleaseLock(lockKey, uniqueValueA)
		assert.NoError(t, err, "Client A trying to release a stale lock should not return an error, but be a no-op/warning")

		// Verify that Client A's attempt did NOT delete Client B's lock
		redisClient := redis.NewClient(&redis.Options{Addr: addr})
		val, err := redisClient.Get(context.Background(), lockKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, uniqueValueB, val, "Client B's lock should still be present and owned by Client B")

		// Client B releases its legitimate lock
		err = locker.ReleaseLock(lockKey, uniqueValueB)
		assert.NoError(t, err, "Client B should successfully release its lock")

		// Verify the key is deleted
		_, err = redisClient.Get(context.Background(), lockKey).Result()
		assert.Equal(t, redis.Nil, err, "Lock key should be deleted after Client B's release")
	})

	t.Run("AcquireLock with zero TTL should still work (no expiration)", func(t *testing.T) {
		lockKey := "test-lock-no-ttl"
		ttl := 0 * time.Second // Zero TTL means no expiration

		acquired, uniqueValue, err := locker.AcquireLock(lockKey, ttl)
		assert.NoError(t, err)
		assert.True(t, acquired, "Should acquire the lock initially with no TTL")
		assert.NotEmpty(t, uniqueValue)

		// Verify it doesn't expire immediately
		time.Sleep(1 * time.Second) // Wait a bit
		redisClient := redis.NewClient(&redis.Options{Addr: addr})
		val, err := redisClient.Get(context.Background(), lockKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, uniqueValue, val, "Lock with zero TTL should not expire")

		// Release it
		err = locker.ReleaseLock(lockKey, uniqueValue)
		assert.NoError(t, err)

		// Verify it's gone
		_, err = redisClient.Get(context.Background(), lockKey).Result()
		assert.Equal(t, redis.Nil, err, "Lock key should be deleted after release")
	})
}
