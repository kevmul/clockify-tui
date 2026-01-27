package models

type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ClientID    string `json:"clientId"`
	ClientName  string `json:"clientName"`
	WorkspaceID string `json:"workspaceId"`
	Color       string `json:"color"`
	IsBillable  bool   `json:"isBillable"`
	IsArchived  bool   `json:"isArchived"`
	Tasks       []Task `json:"tasks,omitempty"`
}
