package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Commands map[string]struct {
		Arguments []string               `yaml:"arguments"`
		Flags     map[string]interface{} `yaml:"flags"`
		Hooks     struct {
			Pre  string `yaml:"pre"`
			Post string `yaml:"post"`
		} `yaml:"hooks"`
	} `yaml:"commands"`
}

func init() {
	viper.SetConfigName("docker-restic")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(os.Getenv("DOCKER_RESTIC_DIR"))
}

func Parse() (Config, error) {
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
