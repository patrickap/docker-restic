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

func BuildCommand(command config.Command) []string {
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

	commandResult := append([]string{command.Binary}, append(command.Arguments, commandFlags...)...)
	return commandResult
}

func ExecuteCommand(args ...string) *exec.Cmd {
	cmd := exec.Command(args[0], args[1:]...)
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
