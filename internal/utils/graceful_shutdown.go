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

var ExitFunc = os.Exit

// signalChan is used only for testing; if nil, a new one is created
var SignalChan chan os.Signal

// gracefulShutdown waits for termination syscalls and doing clean up operations after received it
func GracefullShutdown(ctx context.Context, timeout time.Duration, ops map[string]Operation) <-chan struct{} {
	wait := make(chan struct{})

	go func() {
		s := SignalChan
		if s == nil {
			s = make(chan os.Signal, 1)
			signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		}

		<-s // wait for signal
		log.Println("Shutting down")

		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("Timeout %d ms has been elapsed, force exit", timeout)
			ExitFunc(0)
		})
		defer timeoutFunc.Stop()

		var wg sync.WaitGroup
		for key, op := range ops {
			wg.Add(1)
			go func(k string, fn Operation) {
				defer wg.Done()
				log.Printf("cleaning up %s", k)
				if err := fn(ctx); err != nil {
					log.Printf("%s: clean up failed: %s", k, err.Error())
					return
				}
				log.Printf("%s was shutdown gracefully", k)
			}(key, op)
		}
		wg.Wait()

		close(wait)
	}()

	return wait
}
