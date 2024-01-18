package config

import (
	"fmt"

	"github.com/patrickap/docker-restic/m/v2/internal/env"
	"github.com/patrickap/docker-restic/m/v2/internal/log"
	"github.com/patrickap/docker-restic/m/v2/internal/util/maps"
	"github.com/spf13/viper"
)

type Config struct {
	Repositories map[string]RepositoryConfig `mapstructure:"repositories"`
	Commands     map[string]CommandConfig    `mapstructure:"commands"`
}

type RepositoryConfig struct {
	Options map[string]interface{} `mapstructure:"options"`
}

type CommandConfig struct {
	Arguments []string               `mapstructure:"arguments"`
	Options   map[string]interface{} `mapstructure:"options"`
	Hooks     struct {
		Pre     string `mapstructure:"pre"`
		Post    string `mapstructure:"post"`
		Success string `mapstructure:"success"`
		Failure string `mapstructure:"failure"`
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
		log.Error().Msg("Could not load configuration file")
		panic(configErr)
	}
}

func Current() *Config {
	return config
}

func (c *Config) GetRepositoryList() []string {
	return maps.GetKeys(c.Repositories)
}

func (c *Config) GetCommandList() []string {
	return maps.GetKeys(c.Commands)
}

func (c *CommandConfig) GetOptionList() []string {
	options := []string{}

	for _, option := range maps.SortByKey(c.Options) {
		switch optionType := option.Value.(type) {
		case bool:
			if optionType {
				options = append(options, fmt.Sprintf("--%s", option.Key))
			}
		case string, int:
			options = append(options, fmt.Sprintf("--%s", option.Key), fmt.Sprintf("%v", option.Value))
		case interface{}:
			if optionType, ok := optionType.([]interface{}); ok {
				for _, optionValue := range optionType {
					if optionValue, ok := optionValue.(string); ok {
						options = append(options, fmt.Sprintf("--%s", option.Key), optionValue)
					}
				}
			}
		}
	}

	return options
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
