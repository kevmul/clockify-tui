package projects

import (
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		APIKey:      "test-key",
		WorkspaceId: "workspace-123",
	}

	model := New(cfg)

	if model.config != cfg {
		t.Error("Config should be set correctly")
	}

	if model.ready {
		t.Error("Model should not be ready initially")
	}

	if len(model.projects) != 0 {
		t.Error("Projects should be empty initially")
	}
}

func TestUpdate_WindowSize(t *testing.T) {
	cfg := &config.Config{}
	model := New(cfg)

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ := model.Update(msg)

	// Verify the model was updated (basic check)
	if updatedModel.config != cfg {
		t.Error("Model should maintain config after window size update")
	}
}

func TestUpdate_KeyNavigation(t *testing.T) {
	cfg := &config.Config{}
	model := New(cfg)

	// Set up some test projects
	model.projects = []models.Project{
		{ID: "1", Name: "Project 1"},
		{ID: "2", Name: "Project 2"},
		{ID: "3", Name: "Project 3"},
	}

	// Test down navigation
	downMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedModel, _ := model.Update(downMsg)

	if updatedModel.cursor != 1 {
		t.Errorf("Expected cursor to be 1 after down key, got %d", updatedModel.cursor)
	}

	// Test up navigation
	upMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	updatedModel, _ = updatedModel.Update(upMsg)

	if updatedModel.cursor != 0 {
		t.Errorf("Expected cursor to be 0 after up key, got %d", updatedModel.cursor)
	}

	// Test boundary - up when at top
	updatedModel, _ = updatedModel.Update(upMsg)
	if updatedModel.cursor != 0 {
		t.Error("Cursor should stay at 0 when at top")
	}

	// Test boundary - down when at bottom
	updatedModel.cursor = 2
	downMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedModel, _ = updatedModel.Update(downMsg)
	if updatedModel.cursor != 2 {
		t.Error("Cursor should stay at last position when at bottom")
	}
}

func TestUpdate_ProjectsLoaded(t *testing.T) {
	cfg := &config.Config{}
	model := New(cfg)

	testProjects := []models.Project{
		{ID: "1", Name: "Project 1", ClientName: "Client A"},
		{ID: "2", Name: "Project 2", ClientName: ""},
		{ID: "3", Name: "Project 3", ClientName: "Client B"},
	}

	msg := messages.ProjectsLoadedMsg{Projects: testProjects}
	updatedModel, _ := model.Update(msg)

	if !updatedModel.ready {
		t.Error("Model should be ready after projects loaded")
	}

	if len(updatedModel.projects) != 3 {
		t.Errorf("Expected 3 projects, got %d", len(updatedModel.projects))
	}

	if updatedModel.projects[0].ID != "1" {
		t.Error("Projects should be set correctly")
	}
}

func TestView_NotReady(t *testing.T) {
	cfg := &config.Config{}
	model := New(cfg)

	view := model.View()
	expected := "Loading projects..."

	if view != expected {
		t.Errorf("Expected '%s', got '%s'", expected, view)
	}
}

func TestView_NoProjects(t *testing.T) {
	cfg := &config.Config{}
	model := New(cfg)
	model.ready = true
	model.projects = []models.Project{}

	view := model.View()
	expected := "No projects found."

	if view != expected {
		t.Errorf("Expected '%s', got '%s'", expected, view)
	}
}

func TestView_WithProjects(t *testing.T) {
	cfg := &config.Config{}
	model := New(cfg)
	model.ready = true
	model.projects = []models.Project{
		{ID: "1", Name: "Test Project"},
	}

	// Set up the list items (simulate what happens in Update)
	msg := messages.ProjectsLoadedMsg{Projects: model.projects}
	model, _ = model.Update(msg)

	view := model.View()

	// Should not be the loading or empty message
	if view == "Loading projects..." || view == "No projects found." {
		t.Error("View should render project list when projects are available")
	}
}

func TestItemInterface(t *testing.T) {
	testItem := item{
		title: "Test Project",
		desc:  "project-123",
	}

	if testItem.Title() != "Test Project" {
		t.Error("Title() should return the title")
	}

	if testItem.Description() != "project-123" {
		t.Error("Description() should return the description")
	}

	if testItem.FilterValue() != "Test Project" {
		t.Error("FilterValue() should return the title for filtering")
	}
}
