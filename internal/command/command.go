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
		hookErr := processHook(config.Hooks.Pre...)
		if hookErr != nil {
			return hookErr
		}

		commandErr := lock.RunWithLock(func() error { return processCommand(config.GetCommand()...) })
		if commandErr != nil {
			hookErr := processHook(config.Hooks.Failure...)
			if hookErr != nil {
				return hookErr
			}

			return commandErr
		} else {
			hookErr := processHook(config.Hooks.Success...)
			if hookErr != nil {
				return hookErr
			}
		}

		hookErr = processHook(config.Hooks.Post...)
		if hookErr != nil {
			return hookErr
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
	command.Stdout = &log.LogWrapper{Writer: os.Stdout, Logger: log.Instance().Info}
	command.Stderr = &log.LogWrapper{Writer: os.Stderr, Logger: log.Instance().Error}
	return command.Run()
}
