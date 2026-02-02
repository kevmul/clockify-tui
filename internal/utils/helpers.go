package utils

import (
	"clockify-app/internal/models"
	"fmt"
	"strconv"
	"strings"
	"time"
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

// parseTime converts a time string like "9a" or "3:30p" to a full time.Time
// It handles various formats: 9a, 9:30a, 9, 9:30
func ParseTime(timeStr string, date time.Time) (time.Time, error) {
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
	var n int
	var err error
	if strings.Contains(timeStr, ":") {
		n, err = fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
		if err != nil || n != 2 {
			return time.Time{}, fmt.Errorf("Invalid time format: \"%s\"", timeStr)
		}
	} else {
		n, err = fmt.Sscanf(timeStr, "%d", &hour)
		if err != nil || n != 1 {
			return time.Time{}, fmt.Errorf("Invalid time format: \"%s\"", timeStr)
		}
	}

	// Convert to 24-hour format
	if isPM && hour != 12 {
		hour += 12 // 1pm = 13, 2pm = 14, etc.
	} else if !isPM && hour == 12 {
		hour = 0 // 12am = midnight = 0
	}

	returnTime := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location())

	return returnTime, nil

}
