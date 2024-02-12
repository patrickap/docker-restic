package util

import (
	"bufio"
	"os"
	"os/exec"

	"github.com/patrickap/docker-restic/m/v2/internal/log"
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

	if options.WrapLogs {
		stdout, stdoutErr := cmd.StdoutPipe()
		if stdoutErr != nil {
			return stdoutErr
		}

		stderr, stderrErr := cmd.StderrPipe()
		if stderrErr != nil {
			return stderrErr
		}

		err := cmd.Start()
		if err != nil {
			return err
		}

		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				log.Instance().Info().Msg(scanner.Text())
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				log.Instance().Error().Msg(scanner.Text())
			}
		}()

		return cmd.Wait()
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
}
