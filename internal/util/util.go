package util

import (
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/patrickap/docker-restic/m/v2/internal/config"
)

type Pair struct {
	Key   string
	Value interface{}
}

func CreateCommand(command config.Command) []string {
	commandArgs := command.Arguments
	commandFlags := []string{}

	for _, flag := range SortMapByKey(command.Flags) {
		switch flagType := flag.Value.(type) {
		case bool:
			if flagType {
				commandFlags = append(commandFlags, fmt.Sprintf("--%s", flag.Key))
			}
		default:
			commandFlags = append(commandFlags, fmt.Sprintf("--%s", flag.Key), fmt.Sprintf("%v", flag.Value))
		}
	}

	return append([]string{"restic"}, append(commandArgs, commandFlags...)...)
}

func ExecuteCommand(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
