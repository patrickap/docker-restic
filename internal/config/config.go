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
	commands map[string]ConfigItem `mapstructure:"commands"`
}

type ConfigItem struct {
	command Command `mapstructure:"command"`
	options Options `mapstructure:"options"`
	hooks   Hooks   `mapstructure:"hooks"`
}

type Command []string
type Options map[string]interface{}
type Hooks struct {
	pre     []string `mapstructure:"pre"`
	post    []string `mapstructure:"post"`
	success []string `mapstructure:"success"`
	failure []string `mapstructure:"failure"`
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

func Instance() *Config {
	return config
}

func (c *Config) Commands() map[string]ConfigItem {
	return c.commands
}

func (c *ConfigItem) Command() []string {
	command := append(c.command, c.Options()...)
	return command
}

func (c *ConfigItem) Options() []string {
	options := []string{}

	for _, option := range util.SortByKey(c.options) {
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

func (c *ConfigItem) Hooks() Hooks {
	return c.hooks
}

func (h *Hooks) Pre() []string {
	return h.pre
}

func (h *Hooks) Post() []string {
	return h.post
}

func (h *Hooks) Success() []string {
	return h.success
}

func (h *Hooks) Failure() []string {
	return h.failure
}
