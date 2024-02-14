package util

import (
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

func ExecuteCommand(args ...string) *exec.Cmd {
	return exec.Command(args[0], args[1:]...)
}
