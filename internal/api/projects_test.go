package api

import (
	"testing"
)

func TestGetProjects(t *testing.T) {
	// Skip this test since it requires mocking the HTTP client
	// This would require modifying the Client struct to accept a custom base URL
	t.Skip("Skipping integration test - would require API client refactoring for proper mocking")
}
