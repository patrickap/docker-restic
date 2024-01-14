package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Repository struct {
		Location     string `yaml:"location"`
		PasswordFile string `yaml:"password_file"`
	} `yaml:"repository"`
	Commands map[string]struct {
		Arguments []string               `yaml:"arguments"`
		Flags     map[string]interface{} `yaml:"flags"`
		Hooks     struct {
			Pre  string `yaml:"pre"`
			Post string `yaml:"post"`
		} `yaml:"hooks"`
	} `yaml:"commands"`
}

// TODO: add slog package for better logs + cobra + viper ?!

// TODO: return back error and add parameter for config path + using env
func GetConfig() Config {
	configFile, err := os.ReadFile("/srv/docker-restic/config.yml")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Println("Error parsing config as yaml:", err)
		os.Exit(1)
	}

	return config
}
