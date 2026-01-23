package api

import (
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	debug "clockify-app/internal/utils"
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (c *Client) FetchReports(workspaceId, userId string) ([]models.Report, error) {
	endpoint := "/workspaces/%s/user/%s/reports/summary"
	body, err := c.Get(fmt.Sprintf(endpoint, workspaceId, userId))

	if err != nil {
		return nil, err
	}

	var reports []models.Report
	if err := json.Unmarshal(body, &reports); err != nil {
		return nil, fmt.Errorf("failed to parse reports: %w", err)
	}
	return reports, nil
}

func FetchReports(apiKey, workspaceId, userId string) tea.Cmd {
	return func() tea.Msg {
		client := NewClient(apiKey)
		data, err := client.FetchReports(workspaceId, userId)

		debug.Log("Fetched reports: %d", len(data))

		if err != nil {
			return messages.ErrorMsg{Err: err}
		}

		return messages.ReportsLoadedMsg{
			Reports: data,
		}
	}
}
