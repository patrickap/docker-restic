package command

import (
	"os"

	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/lock"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util"
)

type Runnable struct {
	Run func() error
}

func BuildCommand(config *config.ConfigItem) *Runnable {
	return &Runnable{Run: func() error {
		err := processHook(config.Hooks.Pre...)
		if err != nil {
			return err
		}

		err = lock.RunWithLock(func() error { return processCommand(config.GetCommand()...) })
		if err != nil {
			err := processHook(config.Hooks.Failure...)
			if err != nil {
				return err
			}

			return err
		} else {
			err := processHook(config.Hooks.Success...)
			if err != nil {
				return err
			}
		}

		err = processHook(config.Hooks.Post...)
		if err != nil {
			return err
		}

		return nil
	}}
}

func processHook(args ...string) error {
	if len(args) > 0 {
		hook := util.ExecuteCommand(args...)
		hook.Stdout = os.Stdout
		hook.Stderr = os.Stderr
		return hook.Run()
	}
	return nil
}

func processCommand(args ...string) error {
	command := util.ExecuteCommand(args...)
	command.Stdout = &log.LogWriter{os.Stdout, log.Instance().Info}
	command.Stderr = &log.LogWriter{os.Stderr, log.Instance().Error}
	return command.Run()
}
