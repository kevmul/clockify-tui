package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestEntryJSONSerialization(t *testing.T) {
	entry := Entry{
		ID:          "test-id",
		Description: "Test task",
		ProjectID:   "project-123",
		Duration:    3600, // 1 hour
		TimeInterval: IntervalTime{
			Start: time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		WorkspaceID: "workspace-456",
		UserID:      "user-789",
		Billable:    true,
	}

	// Test marshaling
	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("Failed to marshal entry: %v", err)
	}

	// Test unmarshaling
	var unmarshaled Entry
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal entry: %v", err)
	}

	if unmarshaled.ID != entry.ID {
		t.Errorf("ID mismatch: got %q, want %q", unmarshaled.ID, entry.ID)
	}
	if unmarshaled.Description != entry.Description {
		t.Errorf("Description mismatch: got %q, want %q", unmarshaled.Description, entry.Description)
	}
}

func TestTimeEntryRequestSerialization(t *testing.T) {
	req := TimeEntryRequest{
		Start:       "2024-01-15T09:00:00Z",
		End:         "2024-01-15T10:00:00Z",
		ProjectID:   "project-123",
		Description: "Test task",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	var unmarshaled TimeEntryRequest
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	if unmarshaled.ProjectID != req.ProjectID {
		t.Errorf("ProjectID mismatch: got %q, want %q", unmarshaled.ProjectID, req.ProjectID)
	}
}
