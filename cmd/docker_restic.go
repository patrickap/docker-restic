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

				log.Info().Msgf("Executing hook (pre): %s", command.Hooks.Pre)
				err = util.ExecuteCommand("/bin/sh", "-c", command.Hooks.Pre)
				if err != nil {
					// handle
				}

				commandResult := createResticCommand(command)

				log.Info().Msgf("Executing command: %s", strings.Join(commandResult, " "))
				err = util.ExecuteCommand(commandResult...)
				if err != nil {
					// handle
				}

				log.Info().Msgf("Executing hook (post) %s", command.Hooks.Post)
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

func createResticCommand(command config.Command) []string {
	commandArgs := command.Arguments
	commandFlags := func() []string {
		flags := []string{}
		for _, flag := range util.SortMapByKey(command.Flags) {
			switch flagType := flag.Value.(type) {
			case bool:
				if flagType {
					flags = append(flags, fmt.Sprintf("--%s", flag.Key))
				}
			default:
				flags = append(flags, fmt.Sprintf("--%s", flag.Key), fmt.Sprintf("%v", flag.Value))
			}
		}

		return flags
	}()

	return append([]string{"restic"}, append(commandArgs, commandFlags...)...)
}
