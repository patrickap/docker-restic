package cmd

import (
	"strings"

	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:          "init",
	Short:        "Initialize repositories specified in config file",
	Long:         "Initialize repositories specified in config file: " + strings.Join(util.GetKeys(config.Repositories), ", "),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		for repositoryName := range config.Repositories {
			repositoryConfig := config.Repositories[repositoryName]

			log.Info().Msgf("Initializing repository '%s'", repositoryConfig.Repo)
			// TODO: ExecuteCommand / BuildCommand refactor
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

func init() {
	rootCmd.AddCommand(initCmd)
}
