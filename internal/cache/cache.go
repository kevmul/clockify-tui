package cache

import (
	"clockify-app/internal/models"
	"sync"
	"time"
)

var (
	instance *ClockifyCache
	once     sync.Once
)

type ClockifyCache struct {
	mu sync.RWMutex

	// Cache for entries and projects
	Entries  []models.Entry
	Projects []models.Project

	// Cache for project tasks (loaded on demand)
	ProjectTasks map[string]CachedItem
}

type CachedItem struct {
	Data     interface{}
	CachedAt time.Time
}

func GetInstance() *ClockifyCache {
	once.Do(func() {
		instance = &ClockifyCache{
			ProjectTasks: make(map[string]CachedItem),
		}
	})
	return instance
}

// ================================
// Entries Cache Methods
// ================================

func (c *ClockifyCache) SetEntries(entries []models.Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Entries = entries
}

func (c *ClockifyCache) GetEntries() []models.Entry {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Entries
}

// ================================
// Projects Cache Methods
// ================================

func (c *ClockifyCache) SetProjects(projects []models.Project) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Projects = projects
}

func (c *ClockifyCache) GetProjects() []models.Project {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Projects
}

// ================================
// Project Tasks Cache Methods
// ================================

func (c *ClockifyCache) SetProjectTasks(projectID string, tasks []models.Task) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ProjectTasks[projectID] = CachedItem{
		Data:     tasks,
		CachedAt: time.Now(),
	}
}

func (c *ClockifyCache) GetProjectTasks(projectID string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, exists := c.ProjectTasks[projectID]; exists {
		return item.Data
	}

	return nil
}
