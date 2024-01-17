package cmd

import (
	cfg "github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/spf13/cobra"
)

var config, configErr = cfg.Get()

var rootCmd = &cobra.Command{
	Use:          "docker-restic",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
}

func Execute() error {
	if configErr != nil {
		log.Error().Msg("Could not load configuration file")
		return configErr
	}

	return rootCmd.Execute()
}
