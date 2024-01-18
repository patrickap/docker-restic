package cmd

import (
	"strings"

	"github.com/patrickap/docker-restic/m/v2/internal/command"
	cfg "github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util/maps"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:          "init",
	Short:        "Initialize repositories specified in config file",
	Long:         "Initialize repositories specified in config file: " + strings.Join(maps.GetKeys(config.Repositories), ", "),
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error {
		for repositoryName := range config.Repositories {
			repositoryConfig := config.Repositories[repositoryName]

			log.Info().Msgf("Initializing repository '%s'", repositoryName)
			cmd := command.BuildCommand(cfg.CommandConfig{
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
