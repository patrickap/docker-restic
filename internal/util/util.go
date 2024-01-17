package util

import (
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/anmitsu/go-shlex"
	"github.com/patrickap/docker-restic/m/v2/internal/config"
)

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

// TODO: split commandconfig and hookconfig so i do not have to pass it here
func BuildCommand(config config.CommandConfig) []string {
	binary := "restic"
	arguments := config.Arguments
	flags := []string{}

	if config.Binary != "" {
		binary = config.Binary
	}

	for _, flag := range SortByKey(config.Flags) {
		switch flagType := flag.Value.(type) {
		case bool:
			if flagType {
				flags = append(flags, fmt.Sprintf("--%s", flag.Key))
			}
		default:
			flags = append(flags, fmt.Sprintf("--%s", flag.Key), fmt.Sprintf("%v", flag.Value))
		}
	}

	command := append([]string{binary}, append(arguments, flags...)...)
	return command
}

func ParseCommand(str string) []string {
	strings, err := shlex.Split(str, true)
	if err != nil {
		return []string{}
	}

	return strings
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
