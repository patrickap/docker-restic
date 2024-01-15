package cmd

import (
	"fmt"
	"os"

	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util"
	"github.com/spf13/cobra"
)

// TODO: make only runnable by user restic:restic

var rootCmd = &cobra.Command{
	Use:   "docker-restic",
	Short: "...",
	Long:  "...",
	Run: func(cmd *cobra.Command, args []string) {
		// handle
	},
}

func init() {
	config, err := config.Parse()
	if err != nil {
		// handle
	}

	for commandName := range config.Commands {
		nextCmd := &cobra.Command{
			Use: commandName,
			Run: func(cmd *cobra.Command, args []string) {
				// TODO: create repo folder and init repository if not exists

				commandConfig, exists := config.Commands[commandName]
				if !exists {
					// handle
				}

				err = util.ExecuteCommand("/bin/sh", "-c", commandConfig.Hooks.Pre)
				if err != nil {
					// handle
				}

				// TODO: override flags from config when set on wrapper which take precendence

				resticCmd := append([]string{commandConfig.Command}, commandConfig.Arguments...)
				for key, value := range commandConfig.Flags {
					switch valueType := value.(type) {
					case bool:
						if valueType {
							resticCmd = append(resticCmd, fmt.Sprintf("--%s", key))
						}
					default:
						resticCmd = append(resticCmd, fmt.Sprintf("--%s", key), fmt.Sprintf("%v", value))
					}
				}

				err = util.ExecuteCommand("restic", resticCmd...)
				if err != nil {
					log.Error().Msg("error")
					os.Exit(1)
				}

				err = util.ExecuteCommand("/bin/sh", "-c", commandConfig.Hooks.Post)
				if err != nil {
					// handle
				}
			},
		}

		rootCmd.AddCommand(nextCmd)
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		// handle
		os.Exit(1)
	}
}
