package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string, password string, db int, logger *logrus.Logger) (ports.Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Ping the Redis server to ensure the connection is established.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisLogger := logger.WithFields(logrus.Fields{
		"layer":  "database",
		"driver": addr,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		redisLogger.WithError(err).Error("Redis error")
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{client: client}, nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		logger.WithError(err).Warnf("failed to get key '%s' from Redis: %w", key, err)
		return "", fmt.Errorf("failed to get key '%s' from Redis: %w", key, err)
	}

	return val, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	logger, _ := ctx.Value(logger.LoggerContextKey).(logrus.FieldLogger)

	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		logger.WithError(err).Warnf("failed to set key '%s' in Redis: %w", key, err)
		return fmt.Errorf("failed to set key '%s' in Redis: %w", key, err)
	}
	return nil
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key '%s' from Redis: %w", key, err)
	}
	return nil
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
