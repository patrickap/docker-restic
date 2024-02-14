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
		err := util.ExecuteCommand(&util.ExecuteCommandOptions{Arguments: config.Hooks.Pre})
		if err != nil {
			return err
		}

		command := config.GetCommand()
		err = lock.RunWithLock(func() error {
			return util.ExecuteCommand(&util.ExecuteCommandOptions{Arguments: command, WrapLogs: true})
		})
		if err != nil {
			err := util.ExecuteCommand(&util.ExecuteCommandOptions{Arguments: config.Hooks.Failure})
			if err != nil {
				return err
			}

			return err
		} else {
			err := util.ExecuteCommand(&util.ExecuteCommandOptions{Arguments: config.Hooks.Success})
			if err != nil {
				return err
			}
		}

		err = util.ExecuteCommand(&util.ExecuteCommandOptions{Arguments: config.Hooks.Post})
		if err != nil {
			return err
		}

		return nil
	}}
}
