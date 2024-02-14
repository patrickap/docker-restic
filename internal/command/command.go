package command

import (
	"os"

	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/lock"
	"github.com/patrickap/docker-restic/m/v2/internal/util"
)

type Runnable struct {
	Run func() error
}

func BuildCommand(commandName string, config *config.ConfigItem) *Runnable {
	return &Runnable{Run: func() error {
		if len(config.Hooks.Pre) > 0 {
			pre := util.ExecuteCommand(config.Hooks.Pre...)
			pre.Stdout = os.Stdout
			pre.Stderr = os.Stderr
			err := pre.Run()
			if err != nil {
				return err
			}
		}

		err := lock.RunWithLock(func() error {
			command := util.ExecuteCommand(config.GetCommand()...)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			return command.Run()
		})
		if err != nil {
			if len(config.Hooks.Failure) > 0 {
				failure := util.ExecuteCommand(config.Hooks.Failure...)
				failure.Stdout = os.Stdout
				failure.Stderr = os.Stderr
				err := failure.Run()
				if err != nil {
					return err
				}
			}

			return err
		} else {
			if len(config.Hooks.Success) > 0 {
				success := util.ExecuteCommand(config.Hooks.Success...)
				success.Stdout = os.Stdout
				success.Stderr = os.Stderr
				err := success.Run()
				if err != nil {
					return err
				}
			}
		}

		if len(config.Hooks.Post) > 0 {
			post := util.ExecuteCommand(config.Hooks.Post...)
			post.Stdout = os.Stdout
			post.Stderr = os.Stderr
			err := post.Run()
			if err != nil {
				return err
			}
		}

		return nil
	}}
}
