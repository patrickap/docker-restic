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
		command := config.Command()
		hooks := config.Hooks()

		hookErr := BuildCommandHooks(hooks.Pre()).Run()
		if hookErr != nil {
			return hookErr
		}

		log.Info().Msgf("Running command '%s': %s", commandName, strings.Join(command, " "))
		commandErr := util.ExecuteCommand(command...).Run()
		if commandErr != nil {
			log.Error().Msgf("Failed to run command '%s'", commandName)

			hookErr := BuildCommandHooks(hooks.Failure()).Run()
			if hookErr != nil {
				return hookErr
			}

			return commandErr
		} else {
			hookErr := BuildCommandHooks(hooks.Success()).Run()
			if hookErr != nil {
				return hookErr
			}
		}

		hookErr = BuildCommandHooks(hooks.Post()).Run()
		if hookErr != nil {
			return hookErr
		}

		return nil
	}}
}

func BuildCommandHooks(commandNames []string) *Runnable {
	return &Runnable{Run: func() error {
		for _, commandName := range commandNames {
			commands := config.Instance().Commands()
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
