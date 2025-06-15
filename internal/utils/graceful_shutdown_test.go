package utils_test

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/chud-lori/go-boilerplate/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGracefulShutdown(t *testing.T) {
	// Override ExitFunc to prevent real os.Exit
	calledExit := false
	originalExit := utils.ExitFunc
	utils.ExitFunc = func(code int) {
		calledExit = true
	}
	defer func() { utils.ExitFunc = originalExit }()

	// Inject a mock signal channel
	mockSignal := make(chan os.Signal, 1)
	utils.SignalChan = mockSignal
	defer func() { utils.SignalChan = nil }() // reset after test

	called := false
	ops := map[string]utils.Operation{
		"mock": func(ctx context.Context) error {
			called = true
			return nil
		},
	}

	ctx := context.Background()
	done := utils.GracefullShutdown(ctx, time.Second, ops)

	// Simulate receiving a signal (no actual syscall)
	mockSignal <- syscall.SIGINT

	select {
	case <-done:
		assert.True(t, called, "cleanup should be called")
		assert.False(t, calledExit, "exitFunc should not be called because timeout was not reached")
	case <-time.After(2 * time.Second):
		t.Fatal("test timed out")
	}
}
