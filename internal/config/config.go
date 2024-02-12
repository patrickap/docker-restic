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
	Commands map[string]ConfigItem `mapstructure:"commands"`
}

type ConfigItem struct {
	Command Command `mapstructure:"command"`
	Options Options `mapstructure:"options"`
	Hooks   Hooks   `mapstructure:"hooks"`
}

type Command []string
type Options map[string]interface{}
type Hooks struct {
	Pre     []string `mapstructure:"pre"`
	Post    []string `mapstructure:"post"`
	Success []string `mapstructure:"success"`
	Failure []string `mapstructure:"failure"`
}

var (
	config    *Config
	configErr error
)

func init() {
	viper.SetConfigName("docker-restic")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(env.DOCKER_RESTIC_DIR + "/config")

	config, configErr = parse()
	if configErr != nil {
		log.Log.Error().Msgf("Failed to load config file: %v", configErr)
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

func Instance() *Config {
	return config
}

func (c *Config) GetCommands() map[string]ConfigItem {
	return c.Commands
}

func (c *ConfigItem) GetCommand() []string {
	command := append(c.Command, c.GetOptions()...)
	return command
}

func (c *ConfigItem) GetOptions() []string {
	options := []string{}

	for _, option := range util.GetPairs(c.Options) {
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
