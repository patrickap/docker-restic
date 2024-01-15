package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/patrickap/docker-restic/m/v2/cmd"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		os.Exit(1)
	}()

	cmd.Execute()
}
