package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds backup configuration (name, sources, exclude).
// Extended in later phases with backend, encryption, etc.
type Config struct {
	Name    string   `yaml:"name"`
	Sources []string `yaml:"sources"`
	Exclude []string `yaml:"exclude"`
}

// Load reads and parses a YAML config file from path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}
