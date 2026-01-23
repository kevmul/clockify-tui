package api

import (
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	debug "clockify-app/internal/utils"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// GetEntries fetches time entries for a user in a workspace from Clockify API
func (c *Client) GetEntries(workspaceId, userId string) ([]models.Entry, error) {

	pageSize := "30"
	endpoint := "/workspaces/%s/user/%s/time-entries?page-size=%s"
	body, err := c.Get(fmt.Sprintf(endpoint, workspaceId, userId, pageSize))

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

// CreateTimeEntry creates a new time entry in Clockify
// Takes all the necessary parameters and returns an error if creation fails
func (c *Client) CreateTimeEntry(workspaceID, projectID, description, startTimeStr, endTimeStr string, date time.Time) error {
	// Parse the time range string (e.g., "9a - 5p") into actual times
	startTime := parseTime(startTimeStr, date)
	endTime := parseTime(endTimeStr, date)

	// Build the request payload
	entry := models.TimeEntryRequest{
		Start:       startTime.Format(time.RFC3339), // Convert to RFC3339 format
		End:         endTime.Format(time.RFC3339),
		ProjectID:   projectID,
		Description: strings.ToUpper(description),
	}

	// Build endpoint and make POST request
	endpoint := fmt.Sprintf("/workspaces/%s/time-entries", workspaceID)
	_, err := c.Post(endpoint, entry)

	if err != nil {
		return fmt.Errorf("failed to create time entry: %w", err)
	}

	return nil
}

func (c *Client) UpdateTimeEntry(workspaceID, entryID, projectID, description, startTimeStr, endTimeStr string, date time.Time) error {
	// Parse the time range string (e.g., "9a - 5p") into actual times
	startTime := parseTime(startTimeStr, date)
	endTime := parseTime(endTimeStr, date)

	// Build the request payload
	entry := models.TimeEntryRequest{
		Start:       startTime.Format(time.RFC3339), // Convert to RFC3339 format
		End:         endTime.Format(time.RFC3339),
		ProjectID:   projectID,
		Description: strings.ToUpper(description),
	}

	// Build endpoint and make PUT request
	endpoint := fmt.Sprintf("/workspaces/%s/time-entries/%s", workspaceID, entryID)
	_, err := c.Put(endpoint, entry)

	if err != nil {
		return fmt.Errorf("failed to update time entry: %w", err)
	}

	return nil
}

func (c *Client) DeleteTimeEntry(workspaceID, entryID string) error {
	// Build endpoint and make DELETE request
	debug.Log("Deleting time entry:", entryID)
	endpoint := fmt.Sprintf("/workspaces/%s/time-entries/%s", workspaceID, entryID)
	_, err := c.Delete(endpoint)

	if err != nil {
		return fmt.Errorf("failed to delete time entry: %w", err)
	}

	return nil
}

// parseTime converts a time string like "9a" or "3:30p" to a full time.Time
// It handles various formats: 9a, 9:30a, 9, 9:30
func parseTime(timeStr string, date time.Time) time.Time {
	// Normalize the string: lowercase, remove spaces
	timeStr = strings.ToLower(strings.TrimSpace(timeStr))
	timeStr = strings.ReplaceAll(timeStr, " ", "")

	var hour, minute int

	// Check if PM (afternoon/evening)
	isPM := strings.HasSuffix(timeStr, "p") || strings.HasSuffix(timeStr, "pm")

	// Remove the am/pm suffix
	timeStr = strings.TrimSuffix(strings.TrimSuffix(timeStr, "p"), "m")
	timeStr = strings.TrimSuffix(strings.TrimSuffix(timeStr, "a"), "m")

	// Parse hour and optional minutes
	if strings.Contains(timeStr, ":") {
		fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	} else {
		fmt.Sscanf(timeStr, "%d", &hour)
	}

	// Convert to 24-hour format
	if isPM && hour != 12 {
		hour += 12 // 1pm = 13, 2pm = 14, etc.
	} else if !isPM && hour == 12 {
		hour = 0 // 12am = midnight = 0
	}

	// Combine the date with our parsed time
	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location())
}
