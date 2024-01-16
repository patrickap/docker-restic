package util

import (
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/anmitsu/go-shlex"
	"github.com/patrickap/docker-restic/m/v2/internal/config"
)

type Pair struct {
	Key   string
	Value interface{}
}

func BuildCommand(config config.CommandConfig) []string {
	binary := "restic"
	arguments := config.Arguments
	flags := []string{}

	if config.Binary != "" {
		binary = config.Binary
	}

	for _, flag := range SortMapByKey(config.Flags) {
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

func SortMapByKey(m map[string]interface{}) []Pair {
	pairs := make([]Pair, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, Pair{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Key < pairs[j].Key
	})
	return pairs
}
