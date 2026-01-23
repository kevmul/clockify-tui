package models

// Shared Report model
type Report struct {
	FixedDate    bool   `json:"fixedDate"`
	ID           string `json:"id"`
	IsPublic     bool   `json:"isPublic"`
	Name         string `json:"name"`
	Link         string `json:"link"`
	ReportAuthor string `json:"reportAuthor"`
	Type         string `json:"type"`
}
