package cmd

import (
	"os"
	"strings"

	"github.com/gofrs/flock"
	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "docker-restic",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
}

func init() {
	// TODO: lock file ?!
	// TODO: init repositories ?!
	config, configErr := config.Get()
	if configErr != nil {
		log.Error().Msg("Could not load configuration file")
		os.Exit(1)
	}

	for commandName := range config.Commands {
		childCmd := &cobra.Command{
			Use:          commandName,
			SilenceUsage: true,
			RunE: func(cmd *cobra.Command, args []string) error {
				commandConfig := config.Commands[commandName]

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

		rootCmd.AddCommand(childCmd)
	}
}

func Execute() {
	lock := flock.New(os.Getenv("DOCKER_RESTIC_DIR") + "/docker-restic.lock")
	locked, lockedErr := lock.TryLock()

	if lockedErr != nil {
		log.Error().Msg("Could not acquire lock")
		os.Exit(1)
	}

	if locked {
		commandErr := rootCmd.Execute()
		lock.Unlock()
		if commandErr != nil {
			os.Exit(1)
		}
	}
}
