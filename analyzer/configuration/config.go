package configuration

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		DomainConfig DomainConfig `yaml:"domainConfig"`
	}

	DomainConfig struct {
		Paths []string `yaml:"paths"`
	}
)

func LoadConfig(configFilePath string) (*Config, error) {
	bytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &config, nil
}
