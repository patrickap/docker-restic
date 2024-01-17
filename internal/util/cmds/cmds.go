package cmds

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/anmitsu/go-shlex"
	"github.com/patrickap/docker-restic/m/v2/internal/config"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util/maps"
)

type Command struct {
	Run func() error
}

func BuildCommand(name string, config config.CommandConfig) *Command {
	command := func() []string {
		arguments := config.Arguments
		flags := []string{}

		for _, flag := range maps.SortByKey(config.Flags) {
			switch flagType := flag.Value.(type) {
			case bool:
				if flagType {
					flags = append(flags, fmt.Sprintf("--%s", flag.Key))
				}
			default:
				flags = append(flags, fmt.Sprintf("--%s", flag.Key), fmt.Sprintf("%v", flag.Value))
			}
		}

		command := append([]string{"restic"}, append(arguments, flags...)...)
		return command
	}()

	Run := func() error {
		if config.Hooks.Pre != "" {
			log.Info().Msgf("Executing hook 'pre': %s", config.Hooks.Pre)
			hookErr := ExecuteCommand(ParseCommand(config.Hooks.Pre)...).Run()
			if hookErr != nil {
				log.Error().Msg("Could not execute hook 'pre'")
			}
		}

		log.Info().Msgf("Executing command '%s': %s", name, strings.Join(command, " "))
		commandErr := ExecuteCommand(command...).Run()

		if commandErr != nil {
			log.Error().Msgf("Could not execute command '%s'", name)

			if config.Hooks.Failure != "" {
				log.Info().Msgf("Executing hook 'failure': %s", config.Hooks.Failure)
				hookErr := ExecuteCommand(ParseCommand(config.Hooks.Failure)...).Run()
				if hookErr != nil {
					log.Error().Msg("Could not execute hook 'failure'")
				}
			}
		} else {
			if config.Hooks.Success != "" {
				log.Info().Msgf("Executing hook 'success': %s", config.Hooks.Success)
				hookErr := ExecuteCommand(ParseCommand(config.Hooks.Success)...).Run()
				if hookErr != nil {
					log.Error().Msg("Could not execute hook 'success'")
				}
			}
		}

		if config.Hooks.Post != "" {
			log.Info().Msgf("Executing hook 'post': %s", config.Hooks.Post)
			hookErr := ExecuteCommand(ParseCommand(config.Hooks.Post)...).Run()
			if hookErr != nil {
				log.Error().Msg("Could not execute hook 'post'")
			}
		}

		return commandErr
	}

	return &Command{Run}
}

func ParseCommand(str string) []string {
	strings, err := shlex.Split(str, true)
	if err != nil {
		return []string{}
	}

	return strings
}

func ExecuteCommand(args ...string) *exec.Cmd {
	var cmd *exec.Cmd
	if len(args) > 0 {
		cmd = exec.Command(args[0], args[1:]...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
