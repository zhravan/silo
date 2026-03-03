package config

// Config holds backup configuration (name, sources, exclude).
// Extended in later phases with backend, encryption, etc.
type Config struct {
	Name    string   `yaml:"name"`
	Sources []string `yaml:"sources"`
	Exclude []string `yaml:"exclude"`
}
