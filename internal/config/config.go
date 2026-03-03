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
	return &cfg, nil
}
