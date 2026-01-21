package models

import "time"

// IntervalTime represents a time interval with start and end times and duration
type IntervalTime struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	// Duration int       `json:"duration"` // in seconds
}

type Entry struct {
	ID           string       `json:"id"`
	Description  string       `json:"description"`
	ProjectID    string       `json:"projectId"`
	TaskID       string       `json:"taskId,omitempty"`
	Duration     int          `json:"timeInterval.duration"` // in seconds
	TimeInterval IntervalTime `json:"timeInterval"`
	WorkspaceID  string       `json:"workspaceId"`
	UserID       string       `json:"userId"`
	Billable     bool         `json:"billable"`
	TagIDs       []string     `json:"tagIds,omitempty"`
}
