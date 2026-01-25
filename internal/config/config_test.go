package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigSerialization(t *testing.T) {
	cfg := &Config{
		APIKey:        "test-key",
		UserId:        "user-123",
		WorkspaceId:   "workspace-456",
		WorkspaceName: "Test Workspace",
	}

	// Test JSON marshaling
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	// Test JSON unmarshaling
	var loaded Config
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if loaded.APIKey != cfg.APIKey {
		t.Errorf("APIKey mismatch: got %q, want %q", loaded.APIKey, cfg.APIKey)
	}
	if loaded.WorkspaceId != cfg.WorkspaceId {
		t.Errorf("WorkspaceId mismatch: got %q, want %q", loaded.WorkspaceId, cfg.WorkspaceId)
	}
}

func TestConfigSaveToFile(t *testing.T) {
	// Create temp file
	tmpFile := filepath.Join(t.TempDir(), "test_config.json")

	cfg := &Config{
		APIKey:      "test-key",
		UserId:      "user-123",
		WorkspaceId: "workspace-456",
	}

	// Manually save to temp file
	data, _ := json.MarshalIndent(cfg, "", "  ")
	err := os.WriteFile(tmpFile, data, 0600)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Read back and verify
	readData, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var loaded Config
	err = json.Unmarshal(readData, &loaded)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if loaded.APIKey != cfg.APIKey {
		t.Errorf("APIKey mismatch: got %q, want %q", loaded.APIKey, cfg.APIKey)
	}
}
