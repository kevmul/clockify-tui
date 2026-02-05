package entryform

import (
	"strings"
	"testing"
	"time"

	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"

	tea "github.com/charmbracelet/bubbletea"
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
	if model.step != stepDateSelect {
		t.Errorf("Expected initial step to be %d, got %d", stepDateSelect, model.step)
	}
	if model.editing {
		t.Error("editing should be false initially")
	}
	if model.cursor != 0 {
		t.Error("cursor should be 0 initially")
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
	projects := []models.Project{
		{ID: "proj1", Name: "Project 1"},
	}
	model := New(&config.Config{}, projects)

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
	if !updated.editing {
		t.Error("editing should be true")
	}
	if updated.description.Value() != "Test task" {
		t.Errorf("Expected description 'Test task', got %q", updated.description.Value())
	}
	if updated.selectedProj.ID != "proj1" {
		t.Errorf("Expected selected project ID 'proj1', got %q", updated.selectedProj.ID)
	}
}

func TestFilterProjects(t *testing.T) {
	projects := []models.Project{
		{ID: "proj1", Name: "Web Development"},
		{ID: "proj2", Name: "Mobile App"},
		{ID: "proj3", Name: "Web Design"},
	}
	model := New(&config.Config{}, projects)

	// Test with no search term
	filtered := model.filterProjects()
	if len(filtered) != 3 {
		t.Errorf("Expected 3 projects with no filter, got %d", len(filtered))
	}

	// Test with search term
	model.projectSearch.SetValue("web")
	filtered = model.filterProjects()
	if len(filtered) != 2 {
		t.Errorf("Expected 2 projects matching 'web', got %d", len(filtered))
	}

	// Test case insensitive search
	model.projectSearch.SetValue("WEB")
	filtered = model.filterProjects()
	if len(filtered) != 2 {
		t.Errorf("Expected 2 projects matching 'WEB' (case insensitive), got %d", len(filtered))
	}
}

func TestDateSelectUpdate(t *testing.T) {
	model := New(&config.Config{}, []models.Project{})
	initialDate := model.date

	// Test left arrow (previous day)
	msg := tea.KeyMsg{Type: tea.KeyLeft}
	updated, _ := model.updateDateSelect(msg)
	if !updated.date.Before(initialDate) {
		t.Error("Left arrow should move to previous day")
	}

	// Test right arrow (next day)
	msg = tea.KeyMsg{Type: tea.KeyRight}
	updated, _ = updated.updateDateSelect(msg)
	if !updated.date.Equal(initialDate) {
		t.Error("Right arrow should move to next day")
	}

	// Test 't' key (today)
	model.date = time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
	updated, _ = model.updateDateSelect(msg)
	today := time.Now()
	if !updated.date.Truncate(24 * time.Hour).Equal(today.Truncate(24 * time.Hour)) {
		t.Error("'t' key should set date to today")
	}
}

func TestStepNavigation(t *testing.T) {
	model := New(&config.Config{}, []models.Project{})

	// Test tab navigation forward
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updated, _ := model.Update(msg)
	if updated.step != stepDescriptionInput {
		t.Errorf("Expected step %d after tab, got %d", stepDescriptionInput, updated.step)
	}

	// Test shift+tab navigation backward
	msg = tea.KeyMsg{Type: tea.KeyShiftTab}
	updated, _ = updated.Update(msg)
	if updated.step != stepDateSelect {
		t.Errorf("Expected step %d after shift+tab, got %d", stepDateSelect, updated.step)
	}
}

func TestTasksLoadedMessage(t *testing.T) {
	model := New(&config.Config{}, []models.Project{})
	tasks := []models.Task{
		{ID: "task1", Name: "Task 1"},
		{ID: "task2", Name: "Task 2"},
	}

	msg := messages.TasksLoadedMsg{Tasks: tasks}
	updated, _ := model.Update(msg)

	if len(updated.tasks) != 3 { // 2 tasks + "No Task" option
		t.Errorf("Expected 3 tasks (including 'No Task'), got %d", len(updated.tasks))
	}
	if !updated.tasksReady {
		t.Error("tasksReady should be true after tasks loaded")
	}
	if updated.tasks[2].Name != "No Task" {
		t.Error("Last task should be 'No Task' option")
	}
}

func TestViewMethods(t *testing.T) {
	model := New(&config.Config{}, []models.Project{})

	// Test each view method returns non-empty string
	views := []struct {
		name string
		view func() string
	}{
		{"dateSelect", model.viewDateSelect},
		{"descriptionInput", model.viewDescriptionInput},
		{"projectSelect", model.viewProjectSelect},
		{"timeInput", model.viewTimeInput},
		{"taskInput", model.viewTaskInput},
		{"confirm", model.viewConfirm},
		{"completion", model.viewCompletionInput},
	}

	for _, v := range views {
		result := v.view()
		if strings.TrimSpace(result) == "" {
			t.Errorf("%s view should return non-empty string", v.name)
		}
	}
}

func TestMainView(t *testing.T) {
	model := New(&config.Config{}, []models.Project{})

	// Test view for each step
	for step := stepDateSelect; step <= stepComplete; step++ {
		model.step = step
		view := model.View()
		if strings.TrimSpace(view) == "" {
			t.Errorf("View should return non-empty string for step %d", step)
		}
	}

	// Test unknown step
	model.step = 999
	view := model.View()
	if !strings.Contains(view, "Unknown step") {
		t.Error("View should handle unknown step")
	}
}

func TestGetLines(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"single line", 0},
		{"line1\nline2", 1},
		{"line1\nline2\nline3", 2},
		{"\n\n\n", 3},
	}

	for _, test := range tests {
		result := getLines(test.input)
		if result != test.expected {
			t.Errorf("getLines(%q) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestEscapeKey(t *testing.T) {
	model := New(&config.Config{APIKey: "test", WorkspaceId: "ws1"}, []models.Project{})
	model.step = stepProjectSelect
	model.cursor = 5

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updated, _ := model.Update(msg)

	// Should reset to initial state
	if updated.step != stepDateSelect {
		t.Errorf("Expected step to reset to %d, got %d", stepDateSelect, updated.step)
	}
	if updated.cursor != 0 {
		t.Error("Cursor should reset to 0")
	}
}
