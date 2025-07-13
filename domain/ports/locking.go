package ports

import "time"

type Locking interface {
	AcquireLock(lockKey string, ttl time.Duration) (bool, string, error)
	ReleaseLock(lockKey string, uniqueValue string) error
	Close() error
}
