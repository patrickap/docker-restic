package cmd

import (
	"os"
	"strings"

	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "docker-restic",
	Args: cobra.ExactArgs(1),
}

func init() {
	// TODO: create repo folder and init repository if not exists
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

				// Hook Pre
				log.Info().Msgf("Executing hook 'Pre': %s", commandConfig.Hooks.Pre)
				hookErr := util.ExecuteCommand(strings.Split(commandConfig.Hooks.Pre, " ")...).Run()
				if hookErr != nil {
					log.Error().Msg("Could not execute hook 'Pre'")
				}

				// Command
				command := util.BuildCommand(commandConfig)
				log.Info().Msgf("Executing command '%s': %s", commandName, strings.Join(command, " "))
				commandErr := util.ExecuteCommand(command...).Run()

				if commandErr != nil {
					log.Error().Msgf("Could not execute command '%s'", commandName)

					// Hook Failure
					log.Info().Msgf("Executing hook 'Failure': %s", commandConfig.Hooks.Failure)
					hookErr := util.ExecuteCommand(strings.Split(commandConfig.Hooks.Failure, " ")...).Run()
					if hookErr != nil {
						log.Error().Msg("Could not execute hook 'Failure'")
					}
				} else {
					// Hook Success
					log.Info().Msgf("Executing hook 'Success': %s", commandConfig.Hooks.Success)
					hookErr := util.ExecuteCommand(strings.Split(commandConfig.Hooks.Success, " ")...).Run()
					if hookErr != nil {
						log.Error().Msg("Could not execute hook 'Success'")
					}
				}

				// Hook Post
				log.Info().Msgf("Executing hook 'Post': %s", commandConfig.Hooks.Post)
				hookErr = util.ExecuteCommand(strings.Split(commandConfig.Hooks.Post, " ")...).Run()
				if hookErr != nil {
					log.Error().Msg("Could not execute hook 'Post'")
				}

				return commandErr
			},
		}

		rootCmd.AddCommand(childCmd)
	}
}

func Execute() {
	// TODO: make only runnable by user restic:restic
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
