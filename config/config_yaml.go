package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// LoadYamlConfig loads RootConfig from file.
// The path of the file is specified by environment variable GOMMERCE_CONFIG_PATH,
// If the environment variable is not set, it defaults to "./config/app-deploy.yaml".
func LoadYamlConfig() (RootConfig, error) {
	path, ok := os.LookupEnv("GOMMERCE_CONFIG_PATH")
	if !ok {
		path = "./config/app-deploy.yaml"
	}
	txt, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &rootConfig{}
	if err := yaml.Unmarshal(txt, cfg); err != nil {
		return nil, err
	} else {
		return cfg, nil
	}
}
