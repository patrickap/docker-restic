package cmd

import (
	"strings"

	cfg "github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util/cmds"
	"github.com/patrickap/docker-restic/m/v2/internal/util/maps"
	"github.com/patrickap/docker-restic/m/v2/internal/util/structs"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:          "init",
	Short:        "Initialize repositories specified in config file",
	Long:         "Initialize repositories specified in config file: " + strings.Join(maps.GetKeys(config.Repositories), ", "),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		for repositoryName := range config.Repositories {
			repositoryConfig := config.Repositories[repositoryName]

			arguments := []string{"init"}
			flags, flagsErr := structs.ToMap(repositoryConfig)
			if flagsErr != nil {
				return flagsErr
			}

			log.Info().Msgf("Initializing repository '%s'", repositoryConfig.Repo)
			command := cmds.BuildCommand("init", cfg.CommandConfig{
				Arguments: arguments,
				Flags:     flags,
			})
			commandErr := command.Run()
			if commandErr != nil {
				log.Error().Msgf("Could not initialize repository '%s'", repositoryConfig.Repo)
				return commandErr
			}

			return nil
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
