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

var minTilExpired = 2 * time.Minute

type ClockifyCache struct {
	mu sync.RWMutex

	// Cache for entries and projects
	Entries  CachedItem[[]models.Entry]
	Projects CachedItem[[]models.Project]

	// Cache for project tasks (loaded on demand)
	ProjectTasks map[string]CachedItem[[]models.Task]
}

type CachedItem[T any] struct {
	Data     T
	CachedAt time.Time
}

func GetInstance() *ClockifyCache {
	once.Do(func() {
		instance = &ClockifyCache{
			ProjectTasks: make(map[string]CachedItem[[]models.Task]),
		}
	})
	return instance
}

func (c *ClockifyCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Entries = CachedItem[[]models.Entry]{}
	c.Projects = CachedItem[[]models.Project]{}
	c.ProjectTasks = make(map[string]CachedItem[[]models.Task])
}

// ================================
// Entries Cache Methods
// ================================

func (c *ClockifyCache) SetEntries(entries []models.Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Entries = CachedItem[[]models.Entry]{
		Data:     entries,
		CachedAt: time.Now(),
	}
}

func (c *ClockifyCache) AddEntry(entry models.Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Prepend the new entry into the cached entries
	c.Entries.Data = append([]models.Entry{entry}, c.Entries.Data...)
	c.Entries.CachedAt = time.Now()
}

func (c *ClockifyCache) UpdateEntry(updatedEntry models.Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Find and update the entry in the cached entries
	for i, entry := range c.Entries.Data {
		if entry.ID == updatedEntry.ID {
			c.Entries.Data[i] = updatedEntry
			c.Entries.CachedAt = time.Now()
			return
		}
	}
}

func (c *ClockifyCache) DeleteEntry(entryID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Find and remove the entry from the cached Entries
	for i, entry := range c.Entries.Data {
		if entry.ID == entryID {
			c.Entries.Data = append(c.Entries.Data[:i], c.Entries.Data[i+1:]...)
			c.Entries.CachedAt = time.Now()
			return
		}
	}
}

func (c *ClockifyCache) InvalidateEntries() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Entries = CachedItem[[]models.Entry]{}
}

func (c *ClockifyCache) GetEntries() []models.Entry {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Entries.Data) > 0 {
		if time.Since(c.Entries.CachedAt) < minTilExpired {
			return c.Entries.Data
		}
	}

	return nil
}

// ================================
// Projects Cache Methods
// ================================

func (c *ClockifyCache) SetProjects(projects []models.Project) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Projects = CachedItem[[]models.Project]{
		Data:     projects,
		CachedAt: time.Now(),
	}
}

func (c *ClockifyCache) AddProject(project models.Project) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Append the new project into the cached projects
	c.Projects.Data = append(c.Projects.Data, project)
	c.Projects.CachedAt = time.Now()
}

func (c *ClockifyCache) InvalidateProjects() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Projects = CachedItem[[]models.Project]{}
}

func (c *ClockifyCache) GetProjects() []models.Project {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Projects.Data) > 0 {
		if time.Since(c.Projects.CachedAt) < minTilExpired {
			return c.Projects.Data
		}
	}

	return nil
}

// ================================
// Project Tasks Cache Methods
// ================================

func (c *ClockifyCache) SetProjectTasks(projectID string, tasks []models.Task) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ProjectTasks[projectID] = CachedItem[[]models.Task]{
		Data:     tasks,
		CachedAt: time.Now(),
	}
}

func (c *ClockifyCache) GetProjectTasks(projectID string) []models.Task {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, exists := c.ProjectTasks[projectID]; exists {
		// Check if expired (5 minutes)
		if time.Since(item.CachedAt) < minTilExpired {
			return item.Data
		}
	}

	return nil
}

func (c *ClockifyCache) InvalidateProjectTasks(projectID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.ProjectTasks, projectID)
	delete(c.ProjectTasks, projectID)
}
