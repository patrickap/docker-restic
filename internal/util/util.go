package util

import (
	"bufio"
	"os/exec"

	"github.com/rs/zerolog/log"
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

func ExecuteCommand(args ...string) error {
	if len(args) <= 0 {
		return nil
	}

	cmd := exec.Command(args[0], args[1:]...)

	stdout, stdoutErr := cmd.StdoutPipe()
	if stdoutErr != nil {
		return stdoutErr
	}

	stderr, stderrErr := cmd.StderrPipe()
	if stderrErr != nil {
		return stderrErr
	}

	cmd.Start()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			log.Info().Msg(scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Error().Msg(scanner.Text())
		}
	}()

	return cmd.Wait()

}
