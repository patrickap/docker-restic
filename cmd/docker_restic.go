package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/spf13/cobra"
)

var dockerResticCmd = &cobra.Command{
	Use:   "docker-restic",
	Short: "...",
	Long:  `...`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: init repository if not exists
		resticCmd := args[0]

		config, parseError := config.Parse()
		if parseError != nil {
			// handle
		}

		command, notFound := config.Commands[resticCmd]
		if notFound {
			// handle
		}

		if command.Hooks.Pre != "" {
			preCmd := exec.Command("/bin/sh", "-c", command.Hooks.Pre)
			preCmd.Stdout = os.Stdout
			preCmd.Stderr = os.Stderr
			preCmd.Run()
		}

		// TODO: override flags from config when set on wrapper which take precendence

		resultCmd := exec.Command("restic", append([]string{resticCmd}, command.Arguments...)...)
		for key, value := range command.Flags {
			// TODO: boolean flag parsing handle flag: true -> --flag
			resultCmd.Args = append(resultCmd.Args, fmt.Sprintf("--%s", key), fmt.Sprintf("%v", value))
		}

		resultCmd.Stdout = os.Stdout
		resultCmd.Stderr = os.Stderr

		log.Info().Msg("Running restic command: " + resultCmd.String())
		cmdError := resultCmd.Run()
		if cmdError != nil {
			log.Error().Msg("Error running restic command: " + cmdError.Error())
			os.Exit(1)
		}

		if command.Hooks.Post != "" {
			postCmd := exec.Command("/bin/sh", "-c", command.Hooks.Post)
			postCmd.Stdout = os.Stdout
			postCmd.Stderr = os.Stderr
			postCmd.Run()
		}
	},
}

func Execute() {
	err := dockerResticCmd.Execute()
	if err != nil {
		// handle
		os.Exit(1)
	}
}
