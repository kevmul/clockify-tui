package models

import "time"

type Entry struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	ProjectID   string    `json:"projectId"`
	TaskID      string    `json:"taskId,omitempty"`
	Start       time.Time `json:"timeInterval.start"`
	End         time.Time `json:"timeInterval.end"`
	Duration    int       `json:"timeInterval.duration"` // in seconds
	WorkspaceID string    `json:"workspaceId"`
	UserID      string    `json:"userId"`
	Billable    bool      `json:"billable"`
	TagIDs      []string  `json:"tagIds,omitempty"`
}
