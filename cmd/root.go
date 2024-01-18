package cmd

import (
	cfg "github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/lock"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/spf13/cobra"
)

// TODO: only run as user restic

var config, configErr = cfg.Get()

var rootCmd = &cobra.Command{
	Use:          "docker-restic",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	PersistentPreRunE: func(c *cobra.Command, args []string) error {
		lockErr := lock.Lock()
		if lockErr != nil {
			log.Error().Msg("Could not acquire lock")
			return lockErr
		}

		defer lock.Unlock()

		return nil
	},
}

func Execute() error {
	if configErr != nil {
		log.Error().Msg("Could not load configuration file")
		return configErr
	}

	return rootCmd.Execute()
}
