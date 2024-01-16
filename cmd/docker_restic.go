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

	initCmd := createInitCommand(config)
	childCmds := createChildCommands(config)
	rootCmd.AddCommand(append([]*cobra.Command{initCmd}, childCmds...)...)
}

func Execute() error {
	if configErr != nil {
		log.Error().Msg("Could not load configuration file")
		return configErr
	}

	return rootCmd.Execute()
}

func createInitCommand(config cfg.Config) *cobra.Command {
	initCmd := &cobra.Command{
		Use:          "init",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			for repositoryName := range config.Repositories {
				repositoryConfig := config.Repositories[repositoryName]

				log.Info().Msgf("Initializing repository '%s'", repositoryConfig.Repo)
				commandErr := util.ExecuteCommand("restic", "init", "--repo", repositoryConfig.Repo, "--password-file", repositoryConfig.PasswordFile).Run()
				if commandErr != nil {
					log.Error().Msgf("Could not initialize repository '%s'", repositoryConfig.Repo)
					return commandErr
				}

				return nil
			}

			return nil
		},
	}

	return initCmd
}

func createChildCommands(config cfg.Config) []*cobra.Command {
	childCmds := []*cobra.Command{}

	for commandName := range config.Commands {
		commandConfig := config.Commands[commandName]

		childCmd := &cobra.Command{
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

		childCmds = append(childCmds, childCmd)
	}

	return childCmds
}
