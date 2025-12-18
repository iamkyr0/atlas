package recovery

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func SetupGracefulShutdown(checkpointFunc func() error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		if err := checkpointFunc(); err != nil {
		}
		os.Exit(0)
	}()
}

