package cmd

import (
	"strings"

	"github.com/patrickap/docker-restic/m/v2/internal/lock"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util/cmds"
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
			RunE: func(cmd *cobra.Command, args []string) error {
				log.Info().Msg("Attempting to acquire lock")
				lockErr := lock.Lock()
				if lockErr != nil {
					log.Error().Msg("Could not acquire lock")
					return lockErr
				}
				defer lock.Unlock()

				command := cmds.BuildCommand(commandName, commandConfig)
				commandErr := command.Run()
				return commandErr
			},
		}

		runCmd.AddCommand(runChildCmd)
	}

	rootCmd.AddCommand(runCmd)
}
