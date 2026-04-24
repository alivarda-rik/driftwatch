package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// ServiceConfig represents the declared state of a service.
type ServiceConfig struct {
	Name        string            `yaml:"name" toml:"name"`
	Version     string            `yaml:"version" toml:"version"`
	Image       string            `yaml:"image" toml:"image"`
	Replicas    int               `yaml:"replicas" toml:"replicas"`
	Environment map[string]string `yaml:"environment" toml:"environment"`
	Ports       []string          `yaml:"ports" toml:"ports"`
}

// Load reads a YAML or TOML config file and returns a ServiceConfig.
func Load(path string) (*ServiceConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	var cfg ServiceConfig

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parsing YAML config %q: %w", path, err)
		}
	case ".toml":
		if _, err := toml.Decode(string(data), &cfg); err != nil {
			return nil, fmt.Errorf("parsing TOML config %q: %w", path, err)
		}
	default:
		return nil, fmt.Errorf("unsupported config format %q (use .yaml, .yml, or .toml)", ext)
	}

	if cfg.Name == "" {
		return nil, fmt.Errorf("config %q is missing required field: name", path)
	}

	return &cfg, nil
}
