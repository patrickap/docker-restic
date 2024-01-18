package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "docker-restic",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
}

func Execute() error {
	return rootCmd.Execute()
}
