package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Source struct {
	Key string `yaml:"key"`
	URL string `yaml:"url"`
}

type Config struct {
	Version int      `yaml:"version"`
	Sources []Source `yaml:"sources"`
}

func NewConfig() *Config {
	return &Config{
		Version: 1,
		Sources: []Source{},
	}
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := NewConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}
