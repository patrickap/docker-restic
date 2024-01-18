package shell

import (
	"os"
	"os/exec"

	"github.com/anmitsu/go-shlex"
)

func ParseCommand(command string) []string {
	strings, err := shlex.Split(command, true)
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
