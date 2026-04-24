package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// ServiceConfig represents the declared state of a service.
type ServiceConfig struct {
	Name     string            `yaml:"name"     toml:"name"`
	Version  string            `yaml:"version"  toml:"version"`
	Image    string            `yaml:"image"    toml:"image"`
	Port     int               `yaml:"port"     toml:"port"`
	Replicas int               `yaml:"replicas" toml:"replicas"`
	Env      map[string]string `yaml:"env"      toml:"env"`
}

// Load reads a YAML or TOML file from path and returns a ServiceConfig.
func Load(path string) (*ServiceConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg ServiceConfig
	ext := filepath.Ext(path)

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parsing YAML: %w", err)
		}
	case ".toml":
		if _, err := toml.Decode(string(data), &cfg); err != nil {
			return nil, fmt.Errorf("parsing TOML: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}

	if cfg.Name == "" {
		return nil, fmt.Errorf("config must include a non-empty 'name' field")
	}

	return &cfg, nil
}
