package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the vaultpipe configuration.
type Config struct {
	Vault  VaultConfig  `yaml:"vault"`
	Output OutputConfig `yaml:"output"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address   string            `yaml:"address"`
	Token     string            `yaml:"token"`
	Namespace string            `yaml:"namespace"`
	Secrets   []SecretMapping   `yaml:"secrets"`
}

// SecretMapping maps a Vault path to an env key.
type SecretMapping struct {
	Path string `yaml:"path"`
	Key  string `yaml:"key"`
	Env  string `yaml:"env"`
}

// OutputConfig holds output file settings.
type OutputConfig struct {
	File   string `yaml:"file"`
	Append bool   `yaml:"append"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("vault.address is required")
	}
	if len(c.Vault.Secrets) == 0 {
		return fmt.Errorf("vault.secrets must contain at least one entry")
	}
	for i, s := range c.Vault.Secrets {
		if s.Path == "" {
			return fmt.Errorf("vault.secrets[%d].path is required", i)
		}
		if s.Env == "" {
			return fmt.Errorf("vault.secrets[%d].env is required", i)
		}
	}
	if c.Output.File == "" {
		c.Output.File = ".env"
	}
	return nil
}
