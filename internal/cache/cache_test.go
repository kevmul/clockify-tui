package cache

import (
	"clockify-app/internal/models"
	"strconv"
	"testing"
	"time"
)

func TestGetInstance(t *testing.T) {
	cache1 := GetInstance()
	cache2 := GetInstance()

	if cache1 != cache2 {
		t.Error("GetInstance should return the same instance (singleton)")
	}

	if cache1.ProjectTasks == nil {
		t.Error("ProjectTasks map should be initialized")
	}
}

// Test Entries Cache
func TestEntriesCache(t *testing.T) {
	cache := GetInstance()

	// Test empty cache
	entries := cache.GetEntries()
	if entries != nil {
		t.Error("Expected nil for empty cache")
	}

	// Test setting and getting entries
	testEntries := []models.Entry{
		{ID: "1", Description: "Test entry 1"},
		{ID: "2", Description: "Test entry 2"},
	}

	cache.SetEntries(testEntries)
	retrieved := cache.GetEntries()

	if len(retrieved) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(retrieved))
	}

	if retrieved[0].ID != "1" || retrieved[1].ID != "2" {
		t.Error("Retrieved entries don't match expected values")
	}

	// Test prepending entry with new entry
	testEntry := models.Entry{ID: "3", Description: "New entry"}

	cache.AddEntry(testEntry)
	retrieved = cache.GetEntries()

	if len(retrieved) != 3 {
		t.Errorf("Expected 3 entries after adding, got %d", len(retrieved))
	}

	if retrieved[0].ID != "3" {
		t.Error("New entry was not prepended correctly")
	}

	// Test updating an existing entry
	updatedEntry := models.Entry{ID: "2", Description: "Updated entry 2"}

	cache.UpdateEntry(updatedEntry)
	retrieved = cache.GetEntries()

	if retrieved[2].Description != "Updated entry 2" {
		t.Error("Entry was not updated correctly")
	}

	// Test deleting an entry
	cache.DeleteEntry("1")
	retrieved = cache.GetEntries()

	if len(retrieved) != 2 {
		t.Errorf("Expected 2 entries after deletion, got %d", len(retrieved))
	}

	for _, entry := range retrieved {
		if entry.ID == "1" {
			t.Error("Entry with ID '1' was not deleted")
		}
	}

	// Test invalidating entries
	cache.InvalidateEntries()
	retrieved = cache.GetEntries()
	if retrieved != nil {
		t.Error("Expected nil after invalidating entries")
	}
}

// Test Projects Cache
func TestProjectsCache(t *testing.T) {
	cache := GetInstance()

	// Test empty cache
	projects := cache.GetProjects()
	if projects != nil {
		t.Error("Expected nil for empty cache")
	}

	// Test setting and getting projects
	testProjects := []models.Project{
		{ID: "p1", Name: "Project 1"},
		{ID: "p2", Name: "Project 2"},
	}

	cache.SetProjects(testProjects)
	retrieved := cache.GetProjects()

	if len(retrieved) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(retrieved))
	}

	if retrieved[0].ID != "p1" || retrieved[1].ID != "p2" {
		t.Error("Retrieved projects don't match expected values")
	}

	// Test Adding a project
	newProject := models.Project{ID: "p3", Name: "Project 3"}
	cache.AddProject(newProject)
	retrieved = cache.GetProjects()

	if len(retrieved) != 3 {
		t.Errorf("Expected 3 projects after adding, got %d", len(retrieved))
	}

	if retrieved[2].ID != "p3" {
		t.Error("New project was not added correctly")
	}

	// Test invalidating entries
	cache.InvalidateProjects()
	retrieved = cache.GetProjects()
	if retrieved != nil {
		t.Error("Expected nil after invalidating entries")
	}
}

// Test Project Tasks Cache
func TestProjectTasksCache(t *testing.T) {
	cache := GetInstance()
	projectID := "test-project"

	// Test empty cache
	tasks := cache.GetProjectTasks(projectID)
	if tasks != nil {
		t.Error("Expected nil for empty cache")
	}

	// Test setting and getting project tasks
	testTasks := []models.Task{
		{ID: "t1", Name: "Task 1"},
		{ID: "t2", Name: "Task 2"},
	}

	cache.SetProjectTasks(projectID, testTasks)
	retrieved := cache.GetProjectTasks(projectID)

	if len(retrieved) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(retrieved))
	}

	if retrieved[0].ID != "t1" || retrieved[1].ID != "t2" {
		t.Error("Retrieved tasks don't match expected values")
	}

	// Test invalidating entries
	cache.InvalidateProjectTasks(projectID)
	retrieved = cache.GetProjectTasks(projectID)
	if retrieved != nil {
		t.Error("Expected nil after invalidating entries")
	}
}

// Test Cache Expiration
func TestCacheExpiration(t *testing.T) {
	// Temporarily reduce expiration time for testing
	originalExpiration := minTilExpired
	minTilExpired = 10 * time.Millisecond
	defer func() { minTilExpired = originalExpiration }()

	cache := GetInstance()

	// Set entries
	testEntries := []models.Entry{{ID: "1", Description: "Test"}}
	cache.SetEntries(testEntries)

	// Should get cached data immediately
	retrieved := cache.GetEntries()
	if len(retrieved) != 1 {
		t.Error("Should return cached entries immediately")
	}

	// Wait for expiration
	time.Sleep(15 * time.Millisecond)

	// Should return nil after expiration
	expired := cache.GetEntries()
	if expired != nil {
		t.Error("Should return nil after cache expiration")
	}
}

// Test Concurrent Access
func TestConcurrentAccess(t *testing.T) {
	cache := GetInstance()
	done := make(chan bool, 2)

	// Concurrent writes
	go func() {
		for i := range 100 {
			cache.SetEntries([]models.Entry{{ID: "concurrent" + strconv.Itoa(i)}})
		}
		done <- true
	}()

	go func() {
		for range 100 {
			cache.GetEntries()
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	// Test should complete without race conditions
}
