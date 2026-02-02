package utils

import (
	"testing"
	"time"

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

func TestParseTime(t *testing.T) {
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)

	happyTests := []struct {
		input    string
		expected time.Time
	}{
		{"9a", time.Date(2024, 1, 1, 9, 0, 0, 0, time.Local)},
		{"3:30p", time.Date(2024, 1, 1, 15, 30, 0, 0, time.Local)},
		{"12pm", time.Date(2024, 1, 1, 12, 0, 0, 0, time.Local)},
		{"12am", time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)},
		{"7", time.Date(2024, 1, 1, 7, 0, 0, 0, time.Local)},
		{"11:15", time.Date(2024, 1, 1, 11, 15, 0, 0, time.Local)},
		{"4 PM", time.Date(2024, 1, 1, 16, 0, 0, 0, time.Local)},
	}

	for _, test := range happyTests {
		result, err := ParseTime(test.input, date)
		if err != nil {
			t.Errorf("ParseTime(%q) returned unexpected error: %v", test.input, err)
		}
		if !result.Equal(test.expected) {
			t.Errorf("ParseTime(%q) = %v; want %v", test.input, result, test.expected)
		}
	}

	// Testing if can handle errors
	invalidTests := []struct {
		input    string
		expected string
	}{
		{"abc", "Invalid time format: \"abc\""}, // Invalid input defaults to midnight
	}

	for _, test := range invalidTests {
		_, err := ParseTime(test.input, date)

		if err != nil {
			if err.Error() == test.expected {
				// Expected error. Passed.
				continue
			}
			t.Errorf("ParseTime(%q) returned unexpected error: %v", test.input, err)
		} else {
			t.Errorf("ParseTime(%q) expected error but got none", test.input)
		}
	}
}
