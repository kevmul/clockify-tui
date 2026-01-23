package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey        string `json:"api_key"`
	UserId        string `json:"user_id"`
	WorkspaceId   string `json:"workspace_id"`
	WorkspaceName string `json:"workspace_name"`
}

func LoadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return &Config{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return &Config{}, nil
	}

	var cfg Config
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save writes the configuration to the config file.
func (c *Config) Save() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "clockify-tui", "config.json"), nil
}
