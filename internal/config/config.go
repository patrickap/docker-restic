package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Commands map[string]Command `yaml:"commands"`
}

type Command struct {
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
	viper.AddConfigPath(os.Getenv("DOCKER_RESTIC_DIR"))
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
