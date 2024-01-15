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
			PreRun: func(cmd *cobra.Command, args []string) {
				log.Info().Msg("PreRun called")
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				Execute := func() error {
					command, found := config.Commands[commandName]
					if !found {
						log.Error().Msg("Could not find command in configuration file")
						return errors.New("config not found")
					}

					commandResult := util.BuildCommand(command)

					log.Info().Msgf("Executing restic: %s", strings.Join(commandResult, " "))
					err = util.ExecuteCommand(commandResult...).Run()
					if err != nil {
						log.Error().Msg("Could not execute restic")
						return err
					}

					return nil
				}

				err := Execute()
				if err != nil {
					log.Info().Msg("RunE error called")
					return err
				} else {
					log.Info().Msg("RunE success called")
					return nil
				}
			},
			PostRun: func(cmd *cobra.Command, args []string) {
				log.Info().Msg("PostRun called")
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
