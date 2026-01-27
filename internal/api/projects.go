package api

import (
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// GetProjects fetches all projects for a given workspace
// Returns a slice of Project structs or an error
func (c *Client) GetProjects(workspaceID string) ([]models.Project, error) {
	// Build the endpoint URL with the workspace ID
	pageSize := "1000" // Adjust page size as needed
	endpoint := fmt.Sprintf("/workspaces/%s/projects?page-size=%s&archived=false", workspaceID, pageSize)

	// Make the GET request
	body, err := c.Get(endpoint)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response into a slice of Project structs
	var projects []models.Project
	if err := json.Unmarshal(body, &projects); err != nil {
		return nil, fmt.Errorf("failed to parse projects: %w", err)
	}

	return projects, nil
}

// FetchProjects returns a command that fetches all projects for a given workspace
func FetchProjects(apiKey, workspaceId string) tea.Cmd {
	return func() tea.Msg {
		client := NewClient(apiKey)
		projects, err := client.GetProjects(workspaceId)

		if err != nil {
			return messages.ErrorMsg{Err: err}
		}

		return messages.ProjectsLoadedMsg{
			Projects: projects,
		}
	}
}
