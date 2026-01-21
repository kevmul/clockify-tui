package api

import (
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	debug "clockify-app/internal/utils"
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// GetEntries fetches time entries for a user in a workspace from Clockify API
func (c *Client) GetEntries(workspaceId, userId string) ([]models.Entry, error) {
	// Implementation to fetch entries from Clockify API
	endpoint := "/workspaces/%s/user/%s/time-entries"
	body, err := c.Get(fmt.Sprintf(endpoint, workspaceId, userId))

	if err != nil {
		return nil, err
	}

	var entries []models.Entry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse entries: %w", err)
	}

	return entries, nil
}

// FetchEntries returns a command that fetches time entries for a user in a workspace
func FetchEntries(apiKey, workspaceId, userId string) tea.Cmd {
	debug.Log("Fetching entries...")
	return func() tea.Msg {
		client := NewClient(apiKey)
		entries, err := client.GetEntries(workspaceId, userId)

		if err != nil {
			return messages.ErrorMsg{Err: err}
		}

		return messages.EntriesLoadedMsg{
			Entries: entries,
		}
	}
}
