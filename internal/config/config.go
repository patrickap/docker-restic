package config

import (
	"github.com/patrickap/docker-restic/m/v2/internal/env"
	"github.com/spf13/viper"
)

type Config struct {
	Repositories map[string]RepositoryConfig `yaml:"repositories"`
	Commands     map[string]CommandConfig    `yaml:"commands"`
}

type RepositoryConfig struct {
	Repo         string `yaml:"repo"`
	PasswordFile string `yaml:"password_file"`
}

type CommandConfig struct {
	Binary    string                 `yaml:"binary"`
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
