package api

import (
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (c *Client) GetTasks(workspaceID, projectID string) ([]models.Task, error) {
	// Build the endpoint URL with the workspace ID and project ID
	pageSize := "100" // Adjust page size as needed
	endpoint := fmt.Sprintf("/workspaces/%s/projects/%s/tasks?page-size=%s", workspaceID, projectID, pageSize)

	// Make the GET request
	body, err := c.Get(endpoint)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response into a slice of Task structs
	var tasks []models.Task
	if err := json.Unmarshal(body, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse tasks: %w", err)
	}

	return tasks, nil
}

// FetchTasks returns a command that fetches all tasks for a given project in a workspace
func FetchTasks(apiKey, workspaceId, projectId string) tea.Cmd {
	return func() tea.Msg {
		client := NewClient(apiKey)
		tasks, err := client.GetTasks(workspaceId, projectId)

		if err != nil {
			return messages.ErrorMsg{Err: err}
		}

		return messages.TasksLoadedMsg{
			Tasks: tasks,
		}
	}
}
