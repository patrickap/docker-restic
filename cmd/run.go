package cmd

import (
	"fmt"

	"github.com/patrickap/docker-restic/m/v2/internal/command"
	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/lock"
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
		commandConfig := commandConfig

		runChildCmd := &cobra.Command{
			Use:          commandName,
			SilenceUsage: true,
			RunE: func(c *cobra.Command, args []string) error {
				return lock.RunWithLock(func() error {
					cmd := command.BuildCommand(&commandConfig)
					return cmd.Run()
				})
			},
		}

		runCmd.AddCommand(runChildCmd)
	}

	rootCmd.AddCommand(runCmd)
}
