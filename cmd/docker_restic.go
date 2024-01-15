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
	config, err := config.Parse()
	if err != nil {
		log.Error().Msg("Could not load configuration file")
		os.Exit(1)
	}

	for commandName := range config.Commands {
		subCmd := &cobra.Command{
			Use: commandName,
			Run: func(cmd *cobra.Command, args []string) {
				// TODO: create repo folder and init repository if not exists

				command, exists := config.Commands[commandName]
				if !exists {
					log.Error().Msg("Could not find command in configuration file")
					os.Exit(1)
				}

				log.Info().Msgf("Executing hook (pre): %s", command.Hooks.Pre)
				err = util.ExecuteCommand("/bin/sh", "-c", command.Hooks.Pre)
				if err != nil {
					log.Warn().Msg("Could not execute hook (pre)")
				}

				commandResult := util.CreateCommand(command)

				log.Info().Msgf("Executing restic: %s", strings.Join(commandResult, " "))
				err = util.ExecuteCommand(commandResult...)
				if err != nil {
					log.Error().Msg("Could not execute restic")
					os.Exit(1)
				}

				log.Info().Msgf("Executing hook (post) %s", command.Hooks.Post)
				err = util.ExecuteCommand("/bin/sh", "-c", command.Hooks.Post)
				if err != nil {
					log.Warn().Msg("Could not execute hook (post)")
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
		Cleanup()
		os.Exit(1)
	}
}

func Cleanup() {
	log.Info().Msg("Running cleanup ...")
}
