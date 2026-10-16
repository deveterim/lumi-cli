package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultBaseURL = "https://api.lumi.email/v1"
	ConfigDirName  = ".lumi"
	ConfigFileName = "config.json"
)

type Config struct {
	APIKey      string `json:"api_key"`
	BaseURL     string `json:"base_url"`
	DefaultFrom string `json:"default_from,omitempty"`
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ConfigDirName, ConfigFileName), nil
}

func Load() (*Config, error) {
	cfg := &Config{
		BaseURL: DefaultBaseURL,
	}

	configPath, err := GetConfigPath()
	if err == nil {
		if data, err := os.ReadFile(configPath); err == nil {
			_ = json.Unmarshal(data, cfg)
		}
	}

	// Environment variable overrides (Lumi primary, Lumi fallback)
	if envKey := os.Getenv("LUMI_API_KEY"); envKey != "" {
		cfg.APIKey = envKey
	} else if envKey := os.Getenv("LUMI_API_KEY"); envKey != "" {
		cfg.APIKey = envKey
	}

	if envURL := os.Getenv("LUMI_BASE_URL"); envURL != "" {
		cfg.BaseURL = strings.TrimRight(envURL, "/")
	} else if envURL := os.Getenv("LUMI_BASE_URL"); envURL != "" {
		cfg.BaseURL = strings.TrimRight(envURL, "/")
	}

	if envFrom := os.Getenv("LUMI_DEFAULT_FROM"); envFrom != "" {
		cfg.DefaultFrom = envFrom
	} else if envFrom := os.Getenv("LUMI_DEFAULT_FROM"); envFrom != "" {
		cfg.DefaultFrom = envFrom
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}

	return cfg, nil
}

func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}
