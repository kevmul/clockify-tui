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
// View messages
// =====================================

type SwitchViewMsg struct {
	View int // e.g., 0 = SettingsView, 1 = EntriesView, etc.
}

type ExitViewMsg struct{}

// =====================================
// Data Loading messages
// =====================================

type UserLoadedMsg struct {
	UserId string
}

type ProjectsLoadedMsg struct {
	Projects []models.Project
}

type ProjectSelectedMsg struct {
	Project models.Project
}

type TasksLoadedMsg struct {
	Tasks []models.Task
}

type AllTasksLoadedMsg struct {
	Tasks map[string][]models.Task // map[ProjectID][]Task
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

type EntryDeleteStartedMsg struct {
	EntryId string
}

type EntryDeletedMsg struct {
	EntryId string
}

type EntryUpdateStartedMsg struct {
	Entry models.Entry
}

type EntryUpdatedMsg struct {
	Entry models.Entry
}

// =====================================
// Modal messages
// =====================================

type ModalClosedMsg struct{}

type ShowModalMsg struct {
	ModalType string // "entry", "help", etc
}

type ItemDeletedMsg struct {
	ID   string
	Type string
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
