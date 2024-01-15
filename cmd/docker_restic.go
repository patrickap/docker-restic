package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "docker-restic",
	Args: cobra.ExactArgs(1),
}

func init() {
	config, err := config.Parse()
	if err != nil {
		// handle
	}

	for commandName := range config.Commands {
		subCmd := &cobra.Command{
			Use: commandName,
			Run: func(cmd *cobra.Command, args []string) {
				// TODO: create repo folder and init repository if not exists

				command, exists := config.Commands[commandName]
				if !exists {
					// handle
				}

				err = util.ExecuteCommand("/bin/sh", "-c", command.Hooks.Pre)
				if err != nil {
					// handle
				}

				// TODO: override flags from config when set on wrapper which take precendence

				commandArgs := command.Arguments
				commandFlags := []string{}

				for key, value := range command.Flags {
					switch valueType := value.(type) {
					case bool:
						if valueType {
							commandFlags = append(commandFlags, fmt.Sprintf("--%s", key))
						}
					default:
						commandFlags = append(commandFlags, fmt.Sprintf("--%s", key), fmt.Sprintf("%v", value))
					}
				}

				log.Info().Msg("Running: restic " + strings.Join(append(commandArgs, commandFlags...), " "))

				err = util.ExecuteCommand("restic", append(commandArgs, commandFlags...)...)
				if err != nil {
					log.Error().Msg("error")
					os.Exit(1)
				}

				err = util.ExecuteCommand("/bin/sh", "-c", command.Hooks.Post)
				if err != nil {
					// handle
				}
			},
		}

		rootCmd.AddCommand(subCmd)
	}
}

func Execute() {
	// TODO: make only runnable by user restic:restic

	err := rootCmd.Execute()
	if err != nil {
		// handle
		os.Exit(1)
	}
}
