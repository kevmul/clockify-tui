package messages

import (
	"clockify-app/internal/config"
	"clockify-app/internal/models"
)

// =====================================
// Config messages
// =====================================
type ConfigSavedMsg struct {
	Config      *config.Config
	UserId      string
	WorkspaceId string
}

type ConfigLoadedMsg struct {
	Config *config.Config
}

// =====================================
// Data Loading messages
// =====================================

type ProjectsLoadedMsg struct {
	Projects []models.Project
}

type TasksLoadedMsg struct {
	Tasks []models.Task
}

type EntriesLoadedMsg struct {
	Entries []models.Entry
}

type WorkspacesLoadedMsg struct {
	Workspaces []models.Workspace
}

// =====================================
// Entry messages
// =====================================

type EntrySavedMsg struct {
	Entry models.Entry
}

type EntryDeletedMsg struct {
	EntryId string
}

type EntryUpdatedMsg struct {
	Entries models.Entry
}

// =====================================
// Modal messages
// =====================================

type ModalClosedMsg struct{}

type ShowModalMsg struct {
	ModalType string // "entry", "help", etc
}

// =====================================
// Error messages
// =====================================

type ErrorMsg struct {
	Err error
}

func (e ErrorMsg) Error() string {
	return e.Err.Error()
}
