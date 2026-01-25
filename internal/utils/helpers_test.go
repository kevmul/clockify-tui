package utils

import (
	"testing"

	"clockify-app/internal/models"
)

func TestFindEntryById(t *testing.T) {
	entries := []models.Entry{
		{ID: "1", Description: "Task 1"},
		{ID: "2", Description: "Task 2"},
		{ID: "3", Description: "Task 3"},
	}

	// Test finding existing entry
	entry, err := FindEntryById(entries, "2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if entry.Description != "Task 2" {
		t.Errorf("Expected 'Task 2', got %q", entry.Description)
	}

	// Test finding non-existent entry
	_, err = FindEntryById(entries, "999")
	if err == nil {
		t.Error("Expected error for non-existent entry")
	}
}

func TestFindProjectById(t *testing.T) {
	projects := []models.Project{
		{ID: "1", Name: "Project A"},
		{ID: "2", Name: "Project B"},
		{ID: "3", Name: "Project C"},
	}

	// Test finding existing project
	project, err := FindProjectById(projects, "2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if project.Name != "Project B" {
		t.Errorf("Expected 'Project B', got %q", project.Name)
	}

	// Test finding non-existent project
	_, err = FindProjectById(projects, "999")
	if err == nil {
		t.Error("Expected error for non-existent project")
	}
}
