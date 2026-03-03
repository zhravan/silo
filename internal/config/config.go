package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds full backup configuration.
type Config struct {
	Name        string   `yaml:"name"`
	Sources     []string `yaml:"sources"`
	Exclude     []string `yaml:"exclude"`
	Schedule    Schedule `yaml:"schedule"`
	Compression Compression `yaml:"compression"`
	Encryption  Encryption  `yaml:"encryption"`
	Backend     Backend    `yaml:"backend"`
	Chunking    Chunking   `yaml:"chunking"`
	Concurrency Concurrency `yaml:"concurrency"`
}

// Schedule holds cron and daemon schedule.
type Schedule struct {
	Cron   string     `yaml:"cron"`
	Daemon DaemonSched `yaml:"daemon"`
}

// DaemonSched holds daemon interval and options.
type DaemonSched struct {
	Every          string `yaml:"every"`
	OnStartup      bool   `yaml:"on_startup"`
	JitterSeconds  int    `yaml:"jitter_seconds"`
}

// Compression holds algorithm and level.
type Compression struct {
	Type  string `yaml:"type"`  // zstd, lz4, none
	Level int    `yaml:"level"` // 1-22 for zstd
}

// Encryption holds password env or key file (secrets not in config).
type Encryption struct {
	PasswordEnv string `yaml:"password_env"`
	KeyFile     string `yaml:"key_file"`
}

// Backend holds type and type-specific fields (s3, rest, sftp).
type Backend struct {
	Type string `yaml:"type"` // s3, rest, sftp
	// S3
	Bucket string `yaml:"bucket"`
	Prefix string `yaml:"prefix"`
	Region string `yaml:"region"`
	// REST
	BaseURL    string `yaml:"base_url"`
	AuthHeader string `yaml:"auth_header"`
	TokenEnv   string `yaml:"token_env"`
	// SFTP
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	User    string `yaml:"user"`
	KeyFile string `yaml:"key_file"`
	Path    string `yaml:"path"`
}

// Chunking holds CDC size bounds.
type Chunking struct {
	TargetSize int `yaml:"target_size"`
	MinSize    int `yaml:"min_size"`
	MaxSize    int `yaml:"max_size"`
}

// Concurrency holds worker counts.
type Concurrency struct {
	Upload int `yaml:"upload"`
	Scan   int `yaml:"scan"`
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
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return &cfg, nil
}

// Validate checks config and returns an error on first failure.
func (c *Config) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(c.Sources) == 0 {
		return fmt.Errorf("at least one source is required")
	}
	switch c.Backend.Type {
	case "s3":
		if c.Backend.Bucket == "" {
			return fmt.Errorf("backend.s3: bucket is required")
		}
	case "rest":
		if c.Backend.BaseURL == "" {
			return fmt.Errorf("backend.rest: base_url is required")
		}
	case "sftp":
		if c.Backend.Host == "" || c.Backend.User == "" {
			return fmt.Errorf("backend.sftp: host and user are required")
		}
	case "":
		return fmt.Errorf("backend.type is required (s3, rest, or sftp)")
	default:
		return fmt.Errorf("backend.type must be s3, rest, or sftp (got %q)", c.Backend.Type)
	}
	if c.Encryption.PasswordEnv == "" && c.Encryption.KeyFile == "" {
		return fmt.Errorf("encryption: password_env or key_file is required")
	}
	switch c.Compression.Type {
	case "", "none", "lz4", "zstd":
		// ok
	default:
		return fmt.Errorf("compression.type must be zstd, lz4, or none (got %q)", c.Compression.Type)
	}
	return nil
}
