package cmd

import (
	"strings"

	cfg "github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/lock"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util"
	"github.com/spf13/cobra"
)

var config, configErr = cfg.Get()

var rootCmd = &cobra.Command{
	Use:          "docker-restic",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
}

// TODO: init repositories ?!

func init() {
	if configErr != nil {
		return
	}

	for commandName := range config.Commands {
		commandConfig := config.Commands[commandName]
		childCmd := createChildCommand(commandName, commandConfig)
		rootCmd.AddCommand(childCmd)
	}
}

func Execute() error {
	if configErr != nil {
		log.Error().Msg("Could not load configuration file")
		return configErr
	}

	return rootCmd.Execute()
}

func createChildCommand(name string, config cfg.CommandConfig) *cobra.Command {
	childCmd := &cobra.Command{
		Use:          name,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Msg("Attempting to acquire lock")
			lockErr := lock.Lock()
			if lockErr != nil {
				log.Error().Msg("Could not acquire lock")
				return lockErr
			}
			defer lock.Unlock()

			log.Info().Msgf("Executing hook 'pre': %s", config.Hooks.Pre)
			hookErr := util.ExecuteCommand(util.ParseCommand(config.Hooks.Pre)...).Run()
			if hookErr != nil {
				log.Error().Msg("Could not execute hook 'pre'")
			}

			command := util.BuildCommand(config)
			log.Info().Msgf("Executing command '%s': %s", name, strings.Join(command, " "))
			commandErr := util.ExecuteCommand(command...).Run()

			if commandErr != nil {
				log.Error().Msgf("Could not execute command '%s'", name)

				log.Info().Msgf("Executing hook 'failure': %s", config.Hooks.Failure)
				hookErr := util.ExecuteCommand(util.ParseCommand(config.Hooks.Failure)...).Run()
				if hookErr != nil {
					log.Error().Msg("Could not execute hook 'failure'")
				}
			} else {
				log.Info().Msgf("Executing hook 'success': %s", config.Hooks.Success)
				hookErr := util.ExecuteCommand(util.ParseCommand(config.Hooks.Success)...).Run()
				if hookErr != nil {
					log.Error().Msg("Could not execute hook 'success'")
				}
			}

			log.Info().Msgf("Executing hook 'post': %s", config.Hooks.Post)
			hookErr = util.ExecuteCommand(util.ParseCommand(config.Hooks.Post)...).Run()
			if hookErr != nil {
				log.Error().Msg("Could not execute hook 'post'")
			}

			return commandErr
		},
	}

	return childCmd
}
