package entryform

import (
	"testing"
	"time"

	"clockify-app/internal/config"
	"clockify-app/internal/models"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		APIKey:      "test-key",
		WorkspaceId: "ws1",
	}
	projects := []models.Project{
		{ID: "proj1", Name: "Project 1"},
	}

	model := New(cfg, projects)

	if model.apiKey != cfg.APIKey {
		t.Error("API key not set correctly")
	}
	if model.workspaceID != cfg.WorkspaceId {
		t.Error("Workspace ID not set correctly")
	}
	if len(model.projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(model.projects))
	}
	if model.step != 0 {
		t.Errorf("Expected initial step to be 0, got %d", model.step)
	}
}

func TestSetProjects(t *testing.T) {
	model := New(&config.Config{}, []models.Project{})

	newProjects := []models.Project{
		{ID: "proj1", Name: "Project 1"},
		{ID: "proj2", Name: "Project 2"},
	}

	updated := model.SetProjects(newProjects)

	if len(updated.projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(updated.projects))
	}
}

func TestUpdateEntry(t *testing.T) {
	model := New(&config.Config{}, []models.Project{})

	entry := models.Entry{
		ID:          "entry1",
		Description: "Test task",
		ProjectID:   "proj1",
		TimeInterval: models.IntervalTime{
			Start: time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 15, 17, 0, 0, 0, time.UTC),
		},
	}

	updated := model.UpdateEntry(entry)

	if updated.selectedEntry.ID != "entry1" {
		t.Errorf("Expected entry ID 'entry1', got %q", updated.selectedEntry.ID)
	}
	if updated.editing != true {
		t.Error("editing should be true")
	}
}
