package utils

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// clean up function on shutting down
type Operation func(ctx context.Context) error

// gracefulShutdown waits for termination syscalls and doing clean up operations after received it
func GracefullShutdown(ctx context.Context, timeout time.Duration, ops map[string]Operation) <-chan struct{} {
	wait := make(chan struct{})

	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		log.Println("Shutting down")

		// Set timeout for the ops to be done to prevent hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("Timeout %d ms has been elapsed, force exit", timeout)
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the ops async to save time
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Printf("cleaning up %s", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Printf("%s: clean up failed: %s", innerKey, err.Error())
					return
				}
				log.Printf("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}
