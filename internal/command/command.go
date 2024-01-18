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
		err := BuildCommandHooks(config.Hooks.Pre).Run()
		if err != nil {
			return err
		}

		command := append(config.Command, config.GetOptionList()...)
		log.Info().Msgf("Running command: %s", strings.Join(command, " "))

		err = shell.ExecuteCommand(command...).Run()
		if err != nil {
			err = BuildCommandHooks(config.Hooks.Failure).Run()
			if err != nil {
				return err
			}

			return err
		} else {
			err = BuildCommandHooks(config.Hooks.Success).Run()
			if err != nil {
				return err
			}
		}

		err = BuildCommandHooks(config.Hooks.Post).Run()
		if err != nil {
			return err
		}

		return nil
	}}
}

func BuildCommandHooks(commandNames []string) *Runnable {
	return &Runnable{Run: func() error {
		for _, commandName := range commandNames {
			command, found := config.Current().Commands[commandName]

			if found {
				err := BuildCommand(&command).Run()
				if err != nil {
					return err
				}
			}
		}

		return nil
	}}
}
