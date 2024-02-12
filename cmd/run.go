package cmd

import (
	"github.com/patrickap/docker-restic/m/v2/internal/command"
	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:           "run",
	Short:         "Run provided command specified in config file",
	Long:          "Run provided command specified in config file",
	Args:          cobra.ExactArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	commands := config.Instance().GetCommands()

	for commandName, commandConfig := range commands {
		commandName := commandName
		commandConfig := commandConfig

		runChildCmd := &cobra.Command{
			Use:          commandName,
			SilenceUsage: true,
			RunE: func(c *cobra.Command, args []string) error {
				log.Log.Info().Msgf("Running command: %s", commandName)
				cmd := command.BuildCommand(commandName, &commandConfig)
				cmdErr := cmd.Run()
				if cmdErr != nil {
					log.Log.Error().Msgf("Failed to run command: %s: %v", commandName, cmdErr)
					return cmdErr
				}

				return nil
			},
		}

		runCmd.AddCommand(runChildCmd)
	}

	rootCmd.AddCommand(runCmd)
}
