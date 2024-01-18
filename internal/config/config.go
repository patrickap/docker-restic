package config

import (
	"fmt"

	"github.com/patrickap/docker-restic/m/v2/internal/env"
	"github.com/patrickap/docker-restic/m/v2/internal/util/maps"
	"github.com/spf13/viper"
)

type Config struct {
	Repositories map[string]RepositoryConfig `yaml:"repositories"`
	Commands     map[string]CommandConfig    `yaml:"commands"`
}

type RepositoryConfig map[string]interface{}

type CommandConfig struct {
	Arguments []string               `yaml:"arguments"`
	Flags     map[string]interface{} `yaml:"flags"`
	Hooks     struct {
		Pre     string `yaml:"pre"`
		Post    string `yaml:"post"`
		Success string `yaml:"success"`
		Failure string `yaml:"failure"`
	} `yaml:"hooks"`
}

func init() {
	viper.SetConfigName("docker-restic")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(env.DOCKER_RESTIC_DIR)
}

func Get() (Config, error) {
	var config Config

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
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
