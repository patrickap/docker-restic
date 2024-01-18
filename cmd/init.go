package cmd

import (
	"fmt"

	"github.com/patrickap/docker-restic/m/v2/internal/command"
	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/lock"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:          "init",
	Short:        "Initialize all repositories specified in config file",
	Long:         fmt.Sprintf("Initialize all repositories specified in config file: %v", config.Current().GetRepositoryList()),
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error {
		return lock.RunWithLock(func() error {
			repositories := config.Current().Repositories

			var cmdErr error
			for repositoryName, repositoryConfig := range repositories {
				log.Info().Msgf("Initializing repository '%s'", repositoryName)
				cmd := command.BuildCommand(&config.CommandConfig{
					Arguments: []string{"init"},
					Options:   repositoryConfig.Options,
				})

				cmdErr = cmd.Run()
			}

			return cmdErr
		})
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
