package api

import (
	"clockify-app/internal/models"
	"encoding/json"
)

func (c *Client) GetWorkspaces() ([]models.Workspace, error) {
	var workspaces []models.Workspace
	data, err := c.Get("/workspaces")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &workspaces)
	if err != nil {
		return nil, err
	}

	return workspaces, nil
}
