package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/patrickap/docker-restic/m/v2/internal"
)

func Execute() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: docker-restic [options]")
		os.Exit(1)
	}

	config := internal.GetConfig()

	// TODO: init repository if not exists

	commandName := args[0]
	command, exists := config.Commands[commandName]
	if !exists {
		fmt.Printf("Command '%s' not found in the configuration\n: ", commandName)
		os.Exit(1)
	}

	// TODO: override flags from config when set on wrapper which take precendence

	resticCmd := exec.Command("restic", append([]string{commandName}, command.Arguments...)...)
	for key, value := range command.Flags {
		// TODO: boolean flag parsing handle flag: true -> --flag
		resticCmd.Args = append(resticCmd.Args, fmt.Sprintf("--%s", key), fmt.Sprintf("%v", value))
	}

	if command.Hooks.Pre != "" {
		preCmd := exec.Command("/bin/sh", "-c", command.Hooks.Pre)
		preCmd.Stdout = os.Stdout
		preCmd.Stderr = os.Stderr
		preCmd.Run()
	}

	resticCmd.Stdout = os.Stdout
	resticCmd.Stderr = os.Stderr
	fmt.Println("Running restic command:", resticCmd.String())
	err := resticCmd.Run()
	if err != nil {
		fmt.Println("Error running restic command:", err)
		os.Exit(1)
	}

	if command.Hooks.Post != "" {
		postCmd := exec.Command("/bin/sh", "-c", command.Hooks.Post)
		postCmd.Stdout = os.Stdout
		postCmd.Stderr = os.Stderr
		postCmd.Run()
	}
}
