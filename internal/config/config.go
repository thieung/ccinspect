package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds ccinspect configuration.
type Config struct {
	ScanPaths    []string `json:"scan_paths"`
	ExcludePaths []string `json:"exclude_paths"`
	MaxDepth     int      `json:"max_depth"`
	DefaultOutput string  `json:"default_output"`
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		ScanPaths:    []string{"~/projects", "~/work"},
		ExcludePaths: []string{"node_modules", ".git", "vendor"},
		MaxDepth:     5,
		DefaultOutput: "table",
	}
}

// configPath returns ~/.ccinspect/config.json.
func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ccinspect", "config.json")
}

// Load reads config from disk, falling back to defaults.
func Load() *Config {
	cfg := DefaultConfig()
	data, err := os.ReadFile(configPath())
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(data, cfg) // fallback to defaults on parse error
	return cfg
}

// Save writes config to disk.
func Save(cfg *Config) error {
	p := configPath()
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}
