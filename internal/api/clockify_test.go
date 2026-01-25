package api

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.apiKey != apiKey {
		t.Errorf("Expected API key %q, got %q", apiKey, client.apiKey)
	}

	if client.httpClient == nil {
		t.Error("HTTP client should not be nil")
	}
}
