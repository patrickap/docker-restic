package util

import (
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func SortByKey[K string, V any](m map[K]V) []Pair[K, V] {
	pairs := GetPairs(m)

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Key < pairs[j].Key
	})

	return pairs
}

func GetPairs[K string, V any](m map[K]V) []Pair[K, V] {
	pairs := make([]Pair[K, V], 0, len(m))

	for key, value := range m {
		pairs = append(pairs, Pair[K, V]{key, value})
	}

	return pairs
}

func GetKeys[K string, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	return keys
}

func Replace(s string, replacements map[string]string) string {
	for placeholder, replacement := range replacements {
		if strings.Contains(strings.ToLower(s), strings.ToLower(placeholder)) {
			s = strings.ReplaceAll(s, placeholder, replacement)
		}
	}
	return s
}

func ExecuteCommand(args ...string) *exec.Cmd {
	var cmd *exec.Cmd
	if len(args) > 0 {
		cmd = exec.Command(args[0], args[1:]...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
