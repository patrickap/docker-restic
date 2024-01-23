package util

import (
	"os"
	"os/exec"
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
