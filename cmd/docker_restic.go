package cmd

import (
	"errors"
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
	config, err := config.Get()
	if err != nil {
		log.Error().Msg("Could not load configuration file")
		os.Exit(1)
	}

	for commandName := range config.Commands {
		subCmd := &cobra.Command{
			Use:          commandName,
			SilenceUsage: true,
			RunE: func(cmd *cobra.Command, args []string) error {
				command, found := config.Commands[commandName]
				if !found {
					log.Error().Msg("Could not find command in configuration file")
					return errors.New("config not found")
				}

				log.Info().Msgf("Executing hook 'Pre': %s", command.Hooks.Pre)
				util.ExecuteCommand(strings.Split(command.Hooks.Pre, " ")...).Run()

				commandResult := util.BuildCommand(command)
				log.Info().Msgf("Executing restic: %s", strings.Join(commandResult, " "))
				err = util.ExecuteCommand(commandResult...).Run()

				log.Info().Msgf("Executing hook 'Post': %s", command.Hooks.Post)
				util.ExecuteCommand(strings.Split(command.Hooks.Post, " ")...).Run()

				if err != nil {
					log.Error().Msg("Could not execute restic")

					log.Info().Msgf("Executing hook 'Failure': %s", command.Hooks.Failure)
					util.ExecuteCommand(strings.Split(command.Hooks.Failure, " ")...).Run()

					return err
				} else {
					log.Info().Msgf("Executing hook 'Success': %s", command.Hooks.Success)
					util.ExecuteCommand(strings.Split(command.Hooks.Success, " ")...).Run()

					return nil
				}
			},
		}

		rootCmd.AddCommand(subCmd)
	}
}

func Execute() {
	// TODO: make only runnable by user restic:restic
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
