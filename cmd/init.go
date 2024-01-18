package cmd

import (
	"fmt"

	"github.com/patrickap/docker-restic/m/v2/internal/command"
	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:          "init",
	Short:        "Initialize all repositories specified in config file",
	Long:         fmt.Sprintf("Initialize all repositories specified in config file: %v", config.Current().GetRepositoryList()),
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error {
		repositories := config.Current().Repositories

		for repositoryName, repositoryConfig := range repositories {
			log.Info().Msgf("Initializing repository '%s'", repositoryName)
			cmd := command.BuildCommand(&config.CommandConfig{
				Arguments: []string{"init"},
				Flags:     repositoryConfig,
			})

			cmdErr := cmd.Run()
			if cmdErr != nil {
				log.Error().Msgf("Could not initialize repository '%s'", repositoryName)
				return cmdErr
			}

			return nil
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
