package testutils

import (
	"time"
	"github.com/sony/gobreaker/v2"
)

// FreshBreaker returns a new gobreaker.CircuitBreaker for tests, with a custom name
func FreshBreaker(name string) *gobreaker.CircuitBreaker[[]byte] {
	st := gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     10 * time.Second,
	}
	return gobreaker.NewCircuitBreaker[[]byte](st)
}

// FreshMailBreaker returns a new gobreaker.CircuitBreaker for mail tests
func FreshMailBreaker() *gobreaker.CircuitBreaker[[]byte] {
	st := gobreaker.Settings{
		Name:        "TestApiMailClient",
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     10 * time.Second,
	}
	return gobreaker.NewCircuitBreaker[[]byte](st)
} 