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

func BuildCommand(commandName string, config *config.CommandConfig) *Runnable {
	return &Runnable{Run: func() error {
		preErr := BuildCommandHooks(config.Hooks.Pre).Run()
		if preErr != nil {
			return preErr
		}

		command := append(config.Command, config.GetOptionList()...)
		log.Info().Msgf("Running command '%s': %s", commandName, strings.Join(command, " "))

		commandErr := shell.ExecuteCommand(command...).Run()
		if commandErr != nil {
			log.Error().Msgf("Failed to run command '%s'", commandName)

			failureErr := BuildCommandHooks(config.Hooks.Failure).Run()
			if failureErr != nil {
				return failureErr
			}

			return commandErr
		} else {
			successErr := BuildCommandHooks(config.Hooks.Success).Run()
			if successErr != nil {
				return successErr
			}
		}

		postErr := BuildCommandHooks(config.Hooks.Post).Run()
		if postErr != nil {
			return postErr
		}

		return nil
	}}
}

func BuildCommandHooks(commandNames []string) *Runnable {
	return &Runnable{Run: func() error {
		for _, commandName := range commandNames {
			command, found := config.Current().Commands[commandName]

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
