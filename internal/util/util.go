package util

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/rs/zerolog"
)

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func GetPairs[K string, V any](m map[K]V) []Pair[K, V] {
	pairs := make([]Pair[K, V], 0, len(m))

	for key, value := range m {
		pairs = append(pairs, Pair[K, V]{key, value})
	}

	return pairs
}

type ExecuteCommandOptions struct {
	Arguments []string
	WrapLogs  bool
}

func ExecuteCommand(options *ExecuteCommandOptions) error {
	if len(options.Arguments) <= 0 {
		return nil
	}

	cmd := exec.Command(options.Arguments[0], options.Arguments[1:]...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	var (
		stdout, stderr io.ReadCloser
		err            error
	)

	if options.WrapLogs {
		stdout, err = cmd.StdoutPipe()
		if err != nil {
			return err
		}

		stderr, err = cmd.StderrPipe()
		if err != nil {
			return err
		}
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	if stdout != nil {
		go wrapLogOutput(&wg, log.Instance().Info, stdout)
	}

	if stderr != nil {
		go wrapLogOutput(&wg, log.Instance().Error, stderr)
	}

	wg.Wait()

	return cmd.Wait()
}

func wrapLogOutput(wg *sync.WaitGroup, logger func() *zerolog.Event, reader io.ReadCloser) {
	wg.Add(1)

	defer func() {
		err := reader.Close()
		if err != nil {
			log.Instance().Error().Msgf("Failed to close pipe: %v", err)
		}
		wg.Done()
	}()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		logger().Msg(scanner.Text())
	}

	err := scanner.Err()
	if err != nil {
		log.Instance().Error().Msgf("Failed to process output: %v", err)
	}
}
