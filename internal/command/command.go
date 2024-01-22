package command

import (
	"strings"

	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util"
)

type Runnable struct {
	Run func() error
}

func BuildCommand(commandName string, config *config.ConfigItem) *Runnable {
	return &Runnable{Run: func() error {
		hookErr := BuildCommandHooks(config.Hooks.Pre).Run()
		if hookErr != nil {
			return hookErr
		}

		command := config.GetCommand()
		log.Info().Msgf("Running command '%s': %s", commandName, strings.Join(command, " "))
		commandErr := util.ExecuteCommand(command...).Run()
		if commandErr != nil {
			log.Error().Msgf("Failed to run command '%s'", commandName)

			hookErr := BuildCommandHooks(config.Hooks.Failure).Run()
			if hookErr != nil {
				return hookErr
			}

			return commandErr
		} else {
			hookErr := BuildCommandHooks(config.Hooks.Success).Run()
			if hookErr != nil {
				return hookErr
			}
		}

		hookErr = BuildCommandHooks(config.Hooks.Post).Run()
		if hookErr != nil {
			return hookErr
		}

		return nil
	}}
}

func BuildCommandHooks(commandNames []string) *Runnable {
	return &Runnable{Run: func() error {
		for _, commandName := range commandNames {
			commands := config.Instance().GetCommands()
			command, found := commands[commandName]

			if found {
				err := BuildCommand(commandName, &command).Run()
				if err != nil {
					return err
				}
			}
		}

		return nil
	}}
}
