package config

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/afero"
)

type Source struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

const DefaultSourceName = "working_dir"

type StrategyType string

const (
	StrategyFlatten  StrategyType = "flatten"
	StrategyPreserve StrategyType = "preserve"
	StrategyConcat   StrategyType = "concat"
)

type Include struct {
	Source string
	File   string
}

type Target struct {
	Name          string       `yaml:"name"`
	Output        string       `yaml:"output"`
	StrategyType  StrategyType `yaml:"strategy,omitempty"`
	Include       []string     `yaml:"include"`
	IncludeParsed []Include    `yaml:"-"`
}

type Config struct {
	Version int      `yaml:"version"`
	Sources []Source `yaml:"sources"`
	Targets []Target `yaml:"targets"`
}

func NewConfig() *Config {
	return &Config{
		Version: 1,
		Sources: []Source{},
		Targets: []Target{},
	}
}

// LoadConfig loads configuration from the given path using the provided filesystem.
func LoadConfig(fs afero.Fs, configPath string, workingDir string) (*Config, error) {
	data, err := afero.ReadFile(fs, configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := NewConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := cfg.addWorkingDirSource(workingDir); err != nil {
		return nil, err
	}

	if err := cfg.normalize(); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) addWorkingDirSource(wd string) error {
	hasWorkingDir := false
	for _, source := range c.Sources {
		if source.Name == DefaultSourceName {
			hasWorkingDir = true
			break
		}
	}

	if !hasWorkingDir {
		c.Sources = append([]Source{{Name: DefaultSourceName, URL: wd}}, c.Sources...)
	}

	return nil
}

func (c *Config) normalize() error {
	if err := c.parseIncludes(); err != nil {
		return err
	}
	c.setDefaultSourceForIncludes()
	return nil
}

func (c *Config) parseIncludes() error {
	for i := range c.Targets {
		var includes []Include
		for _, includeStr := range c.Targets[i].Include {
			include, err := ParseInclude(includeStr)
			if err != nil {
				return fmt.Errorf("failed to parse include in target '%s': %w", c.Targets[i].Name, err)
			}
			includes = append(includes, include)
		}
		c.Targets[i].IncludeParsed = includes
	}
	return nil
}

func (c *Config) setDefaultSourceForIncludes() {
	for i := range c.Targets {
		for j := range c.Targets[i].IncludeParsed {
			include := &c.Targets[i].IncludeParsed[j]
			if include.Source == "" {
				include.Source = DefaultSourceName
			}
		}
	}
}

func (c *Config) Validate() error {
	sourceNames := make(map[string]bool)
	for _, source := range c.Sources {
		if source.Name == "" {
			return fmt.Errorf("source name cannot be empty")
		}
		if strings.Contains(source.Name, "/") {
			return fmt.Errorf("source name '%s' cannot contain '/'", source.Name)
		}
		if sourceNames[source.Name] {
			return fmt.Errorf("duplicate source name: %s", source.Name)
		}
		sourceNames[source.Name] = true
	}

	for _, target := range c.Targets {
		if target.StrategyType != "" && target.StrategyType != StrategyFlatten && target.StrategyType != StrategyPreserve && target.StrategyType != StrategyConcat {
			return fmt.Errorf("target '%s' has invalid strategy: %s (must be 'flatten', 'preserve', or 'concat')", target.Name, target.StrategyType)
		}

		for _, include := range target.IncludeParsed {
			if !sourceNames[include.Source] {
				return fmt.Errorf("target '%s' references unknown source: %s", target.Name, include.Source)
			}
		}
	}

	return nil
}

func ParseInclude(includeStr string) (Include, error) {
	// if includeStr starts with @, it's structure is "@source/path"
	if len(includeStr) > 0 && includeStr[0] == '@' {
		parts := strings.SplitN(includeStr[1:], "/", 2)
		if len(parts) != 2 {
			return Include{}, fmt.Errorf("invalid include format: %s", includeStr)
		}

		return Include{
			Source: parts[0],
			File:   parts[1],
		}, nil
	}

	// otherwise, it's just a path in the working_dir source
	return Include{
		Source: DefaultSourceName,
		File:   includeStr,
	}, nil
}
