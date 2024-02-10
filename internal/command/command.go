package command

import (
	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/lock"
	"github.com/patrickap/docker-restic/m/v2/internal/util"
)

type Runnable struct {
	Run func() error
}

func BuildCommand(commandName string, config *config.ConfigItem) *Runnable {
	return &Runnable{Run: func() error {
		hookErr := util.ExecuteCommand(&util.ExecuteCommandOptions{Arguments: config.Hooks.Pre})
		if hookErr != nil {
			return hookErr
		}

		command := config.GetCommand()
		commandErr := lock.RunWithLock(func() error {
			return util.ExecuteCommand(&util.ExecuteCommandOptions{Arguments: command, WrapLogs: true})
		})
		if commandErr != nil {
			hookErr := util.ExecuteCommand(&util.ExecuteCommandOptions{Arguments: config.Hooks.Failure})
			if hookErr != nil {
				return hookErr
			}

			return commandErr
		} else {
			hookErr := util.ExecuteCommand(&util.ExecuteCommandOptions{Arguments: config.Hooks.Success})
			if hookErr != nil {
				return hookErr
			}
		}

		hookErr = util.ExecuteCommand(&util.ExecuteCommandOptions{Arguments: config.Hooks.Post})
		if hookErr != nil {
			return hookErr
		}

		return nil
	}}
}
