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

type RepositoryConfig map[string]interface{}

type CommandConfig struct {
	Arguments []string               `mapstructure:"arguments"`
	Flags     map[string]interface{} `mapstructure:"flags"`
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

func (c *CommandConfig) GetFlagList() []string {
	flags := []string{}

	for _, flag := range maps.SortByKey(c.Flags) {
		switch flagType := flag.Value.(type) {
		case bool:
			if flagType {
				flags = append(flags, fmt.Sprintf("--%s", flag.Key))
			}
		case string, int:
			flags = append(flags, fmt.Sprintf("--%s", flag.Key), fmt.Sprintf("%v", flag.Value))
		case interface{}:
			if flagType, ok := flagType.([]interface{}); ok {
				for _, flagValue := range flagType {
					if flagValue, ok := flagValue.(string); ok {
						flags = append(flags, fmt.Sprintf("--%s", flag.Key), flagValue)
					}
				}
			}
		}
	}

	return flags
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
