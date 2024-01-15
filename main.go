package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/patrickap/docker-restic/m/v2/cmd"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
)

func main() {
	// create channel to receive os signals
	sigs := make(chan os.Signal, 1)

	// register channel to receive sigint and sigterm signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// start goroutine that will perform cleanup when a signal is received
	go func() {
		<-sigs
		cleanup()
		os.Exit(1)
	}()

	cmd.Execute()
}

func cleanup() {
	log.Info().Msg("Running cleanup...")
}
