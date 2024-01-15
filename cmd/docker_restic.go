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
	config, err := config.Get()
	if err != nil {
		log.Error().Msg("Could not load configuration file")
		return
	}

	for commandName := range config.Commands {
		subCmd := &cobra.Command{
			Use: commandName,
			Run: func(cmd *cobra.Command, args []string) {
				// TODO: create repo folder and init repository if not exists

				command, exists := config.Commands[commandName]
				if !exists {
					log.Error().Msg("Could not find command in configuration file")
					return
				}

				defer func() {
					log.Info().Msgf("Executing hook 'post' %s", command.Hooks.Post)
					err = util.ExecuteCommand("/bin/sh", "-c", command.Hooks.Post).Run()
					if err != nil {
						log.Warn().Msg("Could not execute hook 'post'")
					}
				}()

				log.Info().Msgf("Executing hook 'pre': %s", command.Hooks.Pre)
				err = util.ExecuteCommand("/bin/sh", "-c", command.Hooks.Pre).Run()
				if err != nil {
					log.Warn().Msg("Could not execute hook 'pre'")
				}

				commandResult := util.BuildCommand(command)

				log.Info().Msgf("Executing restic: %s", strings.Join(commandResult, " "))
				err = util.ExecuteCommand(commandResult...).Run()
				if err != nil {
					log.Error().Msg("Could not execute restic")

					log.Info().Msgf("Executing hook 'failure' %s", command.Hooks.Failure)
					err = util.ExecuteCommand("/bin/sh", "-c", command.Hooks.Failure).Run()
					if err != nil {
						log.Warn().Msg("Could not execute hook 'failure'")
					}

					return
				}

				log.Info().Msgf("Executing hook 'success' %s", command.Hooks.Success)
				err = util.ExecuteCommand("/bin/sh", "-c", command.Hooks.Success).Run()
				if err != nil {
					log.Warn().Msg("Could not execute hook 'success'")
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
