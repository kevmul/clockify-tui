package utils

import (
	"clockify-app/internal/models"
	"strconv"
)

func FindProjectById(projects []models.Project, id string) (models.Project, error) {
	for _, proj := range projects {
		if proj.ID == id {
			return proj, nil
		}
	}
	return models.Project{}, strconv.ErrSyntax
}

func FindEntryById(entries []models.Entry, id string) (models.Entry, error) {
	for _, entry := range entries {
		if entry.ID == id {
			return entry, nil
		}
	}
	return models.Entry{}, strconv.ErrSyntax
}
