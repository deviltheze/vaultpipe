package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all vaultpipe configuration.
type Config struct {
	VaultAddress string   `yaml:"vault_address"`
	VaultToken   string   `yaml:"vault_token"`
	SecretPath   string   `yaml:"secret_path"`
	OutputFile   string   `yaml:"output_file"`
	TTLSeconds   int      `yaml:"ttl_seconds"`
	MaxBackups   int      `yaml:"max_backups"`
	AuditLog     string   `yaml:"audit_log"`
	Filter       Filter   `yaml:"filter"`
}

// Filter mirrors dotenv.FilterOptions for YAML config.
type Filter struct {
	IncludePrefix []string `yaml:"include_prefix"`
	ExcludePrefix []string `yaml:"exclude_prefix"`
	Keys          []string `yaml:"keys"`
}

// Load reads and parses a YAML config file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.VaultAddress == "" {
		return nil, errors.New("vault_address is required")
	}
	if cfg.OutputFile == "" {
		cfg.OutputFile = ".env"
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = 5
	}
	return &cfg, nil
}
