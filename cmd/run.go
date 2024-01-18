package cmd

import (
	"strings"

	"github.com/patrickap/docker-restic/m/v2/internal/command"
	"github.com/patrickap/docker-restic/m/v2/internal/util/maps"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:          "run",
	Short:        "Run command specified in config file",
	Long:         "Run command specified in config file: " + strings.Join(maps.GetKeys(config.Commands), ", "),
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
}

func init() {
	for commandName := range config.Commands {
		commandConfig := config.Commands[commandName]

		runChildCmd := &cobra.Command{
			Use:          commandName,
			SilenceUsage: true,
			RunE: func(c *cobra.Command, args []string) error {
				cmd := command.BuildCommand(commandConfig)
				cmdErr := cmd.Run()
				return cmdErr
			},
		}

		runCmd.AddCommand(runChildCmd)
	}

	rootCmd.AddCommand(runCmd)
}
