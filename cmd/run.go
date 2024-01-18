package cmd

import (
	"fmt"

	"github.com/patrickap/docker-restic/m/v2/internal/command"
	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:          "run",
	Short:        "Run provided command specified in config file",
	Long:         fmt.Sprintf("Run provided command specified in config file: %v", config.Current().GetCommandList()),
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
}

func init() {
	commands := config.Current().Commands

	for commandName, commandConfig := range commands {
		runChildCmd := &cobra.Command{
			Use:          commandName,
			SilenceUsage: true,
			RunE: func(c *cobra.Command, args []string) error {
				cmd := command.BuildCommand(&commandConfig)
				cmdErr := cmd.Run()
				return cmdErr
			},
		}

		runCmd.AddCommand(runChildCmd)
	}

	rootCmd.AddCommand(runCmd)
}
