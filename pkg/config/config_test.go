package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLoadAndSave(t *testing.T) {
	tempDir := t.TempDir()
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	os.Setenv("HOME", tempDir)
	os.Unsetenv("MAILIFICA_API_KEY")
	os.Unsetenv("MAILIFICA_BASE_URL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}

	if cfg.BaseURL != DefaultBaseURL {
		t.Errorf("expected default base url %s, got %s", DefaultBaseURL, cfg.BaseURL)
	}

	cfg.APIKey = "ma_test_secret_key"
	cfg.BaseURL = "http://localhost:8080/v1"
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Verify file exists
	expectedPath := filepath.Join(tempDir, ConfigDirName, ConfigFileName)
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Fatalf("expected config file at %s", expectedPath)
	}

	// Reload
	reloaded, err := Load()
	if err != nil {
		t.Fatalf("failed to reload config: %v", err)
	}
	if reloaded.APIKey != "ma_test_secret_key" {
		t.Errorf("expected APIKey ma_test_secret_key, got %s", reloaded.APIKey)
	}
	if reloaded.BaseURL != "http://localhost:8080/v1" {
		t.Errorf("expected BaseURL http://localhost:8080/v1, got %s", reloaded.BaseURL)
	}
}

func TestConfigEnvOverrides(t *testing.T) {
	os.Setenv("MAILIFICA_API_KEY", "ma_env_key_123")
	os.Setenv("MAILIFICA_BASE_URL", "https://custom.api.mailifica.com")
	defer os.Unsetenv("MAILIFICA_API_KEY")
	defer os.Unsetenv("MAILIFICA_BASE_URL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.APIKey != "ma_env_key_123" {
		t.Errorf("expected env api key, got %s", cfg.APIKey)
	}
	if cfg.BaseURL != "https://custom.api.mailifica.com" {
		t.Errorf("expected env base url, got %s", cfg.BaseURL)
	}
}
