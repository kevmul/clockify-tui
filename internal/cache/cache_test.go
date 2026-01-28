package cache

import (
	"clockify-app/internal/models"
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
}

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
}

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
}

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

func TestConcurrentAccess(t *testing.T) {
	cache := GetInstance()
	done := make(chan bool, 2)

	// Concurrent writes
	go func() {
		for i := 0; i < 100; i++ {
			cache.SetEntries([]models.Entry{{ID: "concurrent1"}})
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			cache.GetEntries()
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	// Test should complete without race conditions
}
