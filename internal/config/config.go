package config

import (
	"fmt"
	"strings"

	"github.com/patrickap/docker-restic/m/v2/internal/env"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util"
	"github.com/spf13/viper"
)

type Config struct {
	Commands map[string]CommandConfig `mapstructure:"commands"`
}

type CommandConfig struct {
	Command []string               `mapstructure:"command"`
	Options map[string]interface{} `mapstructure:"options"`
	Hooks   struct {
		Pre     []string `mapstructure:"pre"`
		Post    []string `mapstructure:"post"`
		Success []string `mapstructure:"success"`
		Failure []string `mapstructure:"failure"`
	} `mapstructure:"hooks"`
}

var (
	config    *Config
	configErr error
)

func init() {
	viper.SetConfigName("docker-restic")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(env.DOCKER_RESTIC_DIR)

	config, configErr = parse()
	if configErr != nil {
		log.Error().Msg("Failed to load config file")
		panic(configErr)
	}
}

func parse() (*Config, error) {
	var c Config

	err := viper.ReadInConfig()
	if err != nil {
		return &c, err
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return &c, err
	}

	return &c, nil
}

func Current() *Config {
	return config
}

func (c *Config) GetCommandList() []string {
	return util.GetKeys(c.Commands)
}

func (c *CommandConfig) GetCommand(replacements map[string]string) ([]string, bool) {
	replaced := false
	command := make([]string, len(c.Command))
	copy(command, c.Command)

	if replacements != nil {
		for i, word := range command {
			command[i] = util.Replace(word, replacements)
			replaced = true
		}
	}

	return command, replaced
}

func (c *CommandConfig) GetOptionList() []string {
	options := []string{}

	for _, option := range util.SortByKey(c.Options) {
		prefix := "--"
		if strings.HasPrefix(option.Key, "-") {
			prefix = ""
		}

		switch optionType := option.Value.(type) {
		case bool:
			if optionType {
				options = append(options, fmt.Sprintf("%s%s", prefix, option.Key))
			}
		case string, int:
			options = append(options, fmt.Sprintf("%s%s", prefix, option.Key), fmt.Sprintf("%v", option.Value))
		case interface{}:
			if optionType, ok := optionType.([]interface{}); ok {
				for _, optionValue := range optionType {
					if optionValue, ok := optionValue.(string); ok {
						options = append(options, fmt.Sprintf("%s%s", prefix, option.Key), optionValue)
					}
				}
			}
		}
	}

	return options
}
