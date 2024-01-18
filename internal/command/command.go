package command

import (
	"strings"

	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util/shell"
)

type Runnable struct {
	Run func() error
}

func BuildCommand(config *config.CommandConfig) *Runnable {
	return &Runnable{Run: func() error {

		BuildCommandHook("pre", config.Hooks.Pre).Run()

		command := append([]string{"restic"}, append(config.Arguments, config.GetFlagList()...)...)
		log.Info().Msgf("Executing command: %s", strings.Join(command, " "))
		commandErr := shell.ExecuteCommand(command...).Run()

		if commandErr != nil {
			log.Error().Msg("Could not execute command")

			BuildCommandHook("failure", config.Hooks.Failure).Run()
		} else {
			BuildCommandHook("success", config.Hooks.Success).Run()
		}

		BuildCommandHook("post", config.Hooks.Post).Run()

		return commandErr
	}}
}

func BuildCommandHook(name string, command string) *Runnable {
	return &Runnable{Run: func() error {
		if command != "" {
			log.Info().Msgf("Executing hook '%s': %s", name, command)
			hookErr := shell.ExecuteCommand(shell.ParseCommand(command)...).Run()
			if hookErr != nil {
				log.Error().Msgf("Could not execute hook '%s'", name)
				return hookErr
			}
		}

		return nil
	}}
}
