package cmd

import (
	"strings"

	"github.com/patrickap/docker-restic/m/v2/internal/lock"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:          "run",
	Short:        "Run command specified in config file",
	Long:         "Run command specified in config file: " + strings.Join(util.GetKeys(config.Commands), ", "),
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

				log.Info().Msgf("Executing hook 'pre': %s", commandConfig.Hooks.Pre)
				hookErr := util.ExecuteCommand(util.ParseCommand(commandConfig.Hooks.Pre)...).Run()
				if hookErr != nil {
					log.Error().Msg("Could not execute hook 'pre'")
				}

				command := util.BuildCommand(commandConfig)
				log.Info().Msgf("Executing command '%s': %s", commandName, strings.Join(command, " "))
				commandErr := util.ExecuteCommand(command...).Run()

				if commandErr != nil {
					log.Error().Msgf("Could not execute command '%s'", commandName)

					log.Info().Msgf("Executing hook 'failure': %s", commandConfig.Hooks.Failure)
					hookErr := util.ExecuteCommand(util.ParseCommand(commandConfig.Hooks.Failure)...).Run()
					if hookErr != nil {
						log.Error().Msg("Could not execute hook 'failure'")
					}
				} else {
					log.Info().Msgf("Executing hook 'success': %s", commandConfig.Hooks.Success)
					hookErr := util.ExecuteCommand(util.ParseCommand(commandConfig.Hooks.Success)...).Run()
					if hookErr != nil {
						log.Error().Msg("Could not execute hook 'success'")
					}
				}

				log.Info().Msgf("Executing hook 'post': %s", commandConfig.Hooks.Post)
				hookErr = util.ExecuteCommand(util.ParseCommand(commandConfig.Hooks.Post)...).Run()
				if hookErr != nil {
					log.Error().Msg("Could not execute hook 'post'")
				}

				return commandErr
			},
		}

		runCmd.AddCommand(runChildCmd)
	}

	rootCmd.AddCommand(runCmd)
}
