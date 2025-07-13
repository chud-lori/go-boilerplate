package locking

import (
	"context"
	"fmt"
	"time"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// RedisLocker implements the ports.Locking interface using Redis.
type RedisLocker struct {
	client *redis.Client
	logger *logrus.Entry // Use logrus.Entry for context-specific logging
}

// NewRedisLocker creates a new Redis-based locker.
func NewRedisLocker(addr string, password string, db int, logger *logrus.Logger) (ports.Locking, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisLogger := logger.WithFields(logrus.Fields{
		"layer":  "locking_driver", // Different layer for clarity
		"driver": addr,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		redisLogger.WithError(err).Error("Redis locking connection error")
		return nil, fmt.Errorf("failed to connect to Redis for locking: %w", err)
	}

	return &RedisLocker{client: client, logger: redisLogger}, nil
}

// AcquireLock implements the AcquireLock method of the Locking interface.
// It returns true if acquired, the unique value used, and an error.
func (rl *RedisLocker) AcquireLock(lockKey string, ttl time.Duration) (bool, string, error) {
	uniqueValue := uuid.New().String() // Generate a new unique ID for this lock attempt

	// SET key value NX EX seconds
	acquired, err := rl.client.SetNX(context.Background(), lockKey, uniqueValue, ttl).Result()
	if err != nil {
		rl.logger.WithError(err).Errorf("Failed to acquire lock for key: %s", lockKey)
		return false, "", fmt.Errorf("failed to acquire lock: %w", err)
	}

	if !acquired {
		rl.logger.Debugf("Lock '%s' not acquired (already held)", lockKey)
		return false, "", nil
	}

	rl.logger.Debugf("Lock '%s' acquired with value '%s'", lockKey, uniqueValue)

	return acquired, uniqueValue, nil
}

// ReleaseLock implements the ReleaseLock method of the Locking interface.
// It attempts to release a pessimistic lock identified by lockKey, but only if
// the uniqueValue provided matches the value currently stored in Redis for that lock.
//
// This method uses a Lua script to ensure the check (GET) and delete (DEL) operations
// are performed atomically on the Redis server. This is critical to prevent a "stale lock"
// problem, where a client might inadvertently delete a lock that has expired and
// subsequently been acquired by another client.
//
// Parameters:
//   - lockKey: The key representing the lock to be released.
//   - uniqueValue: The unique identifier that was returned by AcquireLock when this
//     client successfully acquired the lock.
//
// Returns:
//   - error: An error if a communication or Redis-specific issue occurred during
//     script execution. It does NOT return an error if the lock was not owned
//     by this client or had already expired (this scenario results in a warning log).
func (rl *RedisLocker) ReleaseLock(lockKey string, uniqueValue string) error {
	// Lua script for atomic lock release.
	// This script ensures that the lock is only deleted if its current value
	// in Redis matches the uniqueValue provided by the client attempting to release it.
	//
	// KEYS[1] will be the lockKey.
	// ARGV[1] will be the uniqueValue (the "signature" of the lock owner).
	script := redis.NewScript(`
		-- Get the current value of the lock key.
		-- If the key does not exist, or its value does not match ARGV[1],
		-- then this client is not the current owner of the lock.
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			-- If the client is the owner, atomically delete the key.
			-- redis.call("DEL", KEYS[1]) returns 1 if the key was deleted, 0 if it didn't exist.
			return redis.call("DEL", KEYS[1])
		else
			-- If the client is not the owner (or the lock expired/was taken by another),
			-- do not delete the key and return 0.
			return 0
		end
	`)

	// Run the Lua script on the Redis server.
	// The `Run` method passes `lockKey` as KEYS[1] and `uniqueValue` as ARGV[1].
	// Redis guarantees that this entire script executes atomically, preventing
	// race conditions between the GET and DEL operations.
	result, err := script.Run(context.Background(), rl.client, []string{lockKey}, uniqueValue).Result()
	if err != nil {
		rl.logger.WithError(err).Errorf("Failed to run Lua script for lock release: %s", lockKey)
		return fmt.Errorf("failed to release lock (script error): %w", err)
	}

	// Interpret the result from the Lua script.
	// The script returns 1 if the lock was successfully deleted, 0 otherwise.
	if res, ok := result.(int64); ok {
		if res == 0 {
			// This scenario means the `DEL` command inside the script was not executed.
			// This typically happens because:
			// 1. The lock key had already expired from Redis.
			// 2. Another client acquired the lock after this client's lock expired,
			//    and thus the unique value no longer matched.
			// This is a normal and expected outcome for distributed locks; it means
			// the client tried to release a lock it no longer owned, and the safety
			// mechanism prevented an incorrect deletion. We log a warning.
			rl.logger.Warnf("Attempted to release lock '%s' with value '%s' but it was not owned or already expired.", lockKey, uniqueValue)
			// Depending on strictness, you *could* return an error here, but typically
			// a non-owned release is just a no-op from a caller's perspective.
		} else {
			// res == 1: The lock was successfully deleted by this client.
			rl.logger.Debugf("Lock '%s' successfully released with value '%s'", lockKey, uniqueValue)
		}
	} else {
		// This block should ideally not be reached if the script consistently returns 0 or 1.
		// It indicates an unexpected return type from the Redis Lua script.
		rl.logger.Errorf("Unexpected result type from Lua script for lock '%s': %T, value: %v", lockKey, result, result)
		return fmt.Errorf("unexpected result type from lock release script")
	}

	return nil
}

func (rl *RedisLocker) Close() error {
	if rl.client == nil {
		return nil // Already closed or not initialized
	}
	rl.logger.Info("Closing Redis locker client connection...")
	err := rl.client.Close()
	if err != nil {
		rl.logger.WithError(err).Error("Failed to close Redis locker client.")
		return fmt.Errorf("failed to close Redis locker client: %w", err)
	}
	rl.logger.Info("Redis locker client connection closed.")
	return nil
}
