package api

import (
	"clockify-app/internal/cache"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"clockify-app/internal/utils"
	"encoding/json"
	"fmt"
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
		// Return cached entries if available
		cache := cache.GetInstance()
		if cachedEntries := cache.GetEntries(); cachedEntries != nil {
			return messages.EntriesLoadedMsg{
				Entries: cachedEntries,
			}
		}

		client := NewClient(apiKey)
		entries, err := client.GetEntries(workspaceId, userId)

		if err != nil {
			return messages.ErrorMsg{Err: err}
		}

		cache.SetEntries(entries)
		return messages.EntriesLoadedMsg{
			Entries: entries,
		}
	}
}

// FetchEntriesForWeek returns a command that fetches time entries for a specific week
func FetchEntriesForWeek(apiKey, workspaceId, userId string, weekStart time.Time) tea.Cmd {
	return func() tea.Msg {
		client := NewClient(apiKey)

		// Calculate start and end of the week
		weekEnd := weekStart.AddDate(0, 0, 7)
		startStr := weekStart.Format("2006-01-02")
		endStr := weekEnd.Format("2006-01-02")

		endpoint := "/workspaces/%s/user/%s/time-entries?start=%sT00:00:00Z&end=%sT00:00:00Z"
		body, err := client.Get(fmt.Sprintf(endpoint, workspaceId, userId, startStr, endStr))

		if err != nil {
			return messages.ErrorMsg{Err: err}
		}

		var entries []models.Entry
		if err := json.Unmarshal(body, &entries); err != nil {
			return messages.ErrorMsg{Err: fmt.Errorf("failed to parse entries: %w", err)}
		}

		return messages.EntriesLoadedMsg{
			Entries: entries,
		}
	}
}

// CreateTimeEntry creates a new time entry in Clockify
// Takes all the necessary parameters and returns an error if creation fails
func (c *Client) CreateTimeEntry(workspaceID, projectID, taskID, description, startTimeStr, endTimeStr string, date time.Time) (models.Entry, error) {

	// Parse the time range string (e.g., "9a - 5p") into actual times
	startTime, _ := utils.ParseTime(startTimeStr, date)
	endTime, _ := utils.ParseTime(endTimeStr, date)

	// Build the request payload
	entry := models.TimeEntryRequest{
		Start:       startTime.Format(time.RFC3339), // Convert to RFC3339 format
		End:         endTime.Format(time.RFC3339),
		ProjectID:   projectID,
		TaskID:      taskID,
		Description: description,
	}

	// Build endpoint and make POST request
	endpoint := fmt.Sprintf("/workspaces/%s/time-entries", workspaceID)
	bytes, err := c.Post(endpoint, entry)

	if err != nil {
		return models.Entry{}, fmt.Errorf("failed to create time entry: %w", err)
	}

	// Parse response
	var newEntry models.Entry
	if err := json.Unmarshal(bytes, &newEntry); err != nil {
		return models.Entry{}, fmt.Errorf("failed to parse created time entry: %w", err)
	}

	return newEntry, nil
}

func (c *Client) UpdateTimeEntry(workspaceID, entryID, projectID, taskID, description, startTimeStr, endTimeStr string, date time.Time) (models.Entry, error) {
	// Parse the time range string (e.g., "9a - 5p") into actual times
	startTime, _ := utils.ParseTime(startTimeStr, date)
	endTime, _ := utils.ParseTime(endTimeStr, date)

	// Build the request payload
	entry := models.TimeEntryRequest{
		Start:       startTime.Format(time.RFC3339), // Convert to RFC3339 format
		End:         endTime.Format(time.RFC3339),
		ProjectID:   projectID,
		TaskID:      taskID,
		Description: description,
	}

	// Build endpoint and make PUT request
	endpoint := fmt.Sprintf("/workspaces/%s/time-entries/%s", workspaceID, entryID)
	bytes, err := c.Put(endpoint, entry)

	if err != nil {
		return models.Entry{}, fmt.Errorf("failed to update time entry: %w", err)
	}

	var updatedEntry models.Entry
	if err := json.Unmarshal(bytes, &updatedEntry); err != nil {
		return models.Entry{}, fmt.Errorf("failed to parse updated time entry: %w", err)
	}

	return updatedEntry, nil
}

func (c *Client) DeleteTimeEntry(workspaceID, entryID string) error {
	// Build endpoint and make DELETE request
	endpoint := fmt.Sprintf("/workspaces/%s/time-entries/%s", workspaceID, entryID)
	_, err := c.Delete(endpoint)

	if err != nil {
		return fmt.Errorf("failed to delete time entry: %w", err)
	}

	return nil
}
