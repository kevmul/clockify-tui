package api

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		input    string
		expected time.Time
	}{
		{"9a", time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)},
		{"9:30a", time.Date(2024, 1, 15, 9, 30, 0, 0, time.UTC)},
		{"12p", time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)},
		{"2p", time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC)},
		{"2:30p", time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)},
		{"12a", time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		{"9", time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)},
		{"14:30", time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)},
		// Edge cases
		{"11:59p", time.Date(2024, 1, 15, 23, 59, 0, 0, time.UTC)},
		{"12:01a", time.Date(2024, 1, 15, 0, 1, 0, 0, time.UTC)},
		{"23:59", time.Date(2024, 1, 15, 23, 59, 0, 0, time.UTC)},
		{"00:00", time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		result := parseTime(tt.input, date)
		if !result.Equal(tt.expected) {
			t.Errorf("parseTime(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestCreateTimeEntry(t *testing.T) {
	// Skip this test since it requires mocking the HTTP client
	t.Skip("Skipping integration test - would require API client refactoring for proper mocking")
}
