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
