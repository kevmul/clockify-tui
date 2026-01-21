package ui

import (
	"clockify-app/internal/api"
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	debug "clockify-app/internal/utils"
	// "clockify-app/internal/ui/views/reports"
	"clockify-app/internal/ui/components/modal"
	"clockify-app/internal/ui/views/entries"
	"clockify-app/internal/ui/views/settings"

	tea "github.com/charmbracelet/bubbletea"
)

type View int

const (
	SettingsView View = iota
	EntriesView
	ReportsView
)

type Model struct {
	// Config and shared state
	config      *config.Config
	userId      string
	workspaceId string
	projects    []models.Project
	tasks       []models.Task

	// Current View
	currentView View

	// View models
	settings settings.Model
	entries  entries.Model
	// reports  reports.Model

	// Modal state
	modal     *modal.Model
	showModal bool

	// UI Dimensions
	width  int
	height int

	// Loading state
	ready bool
}

func NewModel() Model {
	cfg, _ := config.LoadConfig()

	// Start at settings if no config
	currentView := SettingsView

	if cfg.APIKey != "" && cfg.WorkspaceId != "" {
		currentView = EntriesView
	}

	return Model{
		config:      cfg,
		currentView: currentView,
		settings:    settings.New(cfg),
		entries:     entries.New(cfg),
		// reports:     reports.New(),
		ready: false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.initializeFirstViewCmd(),
		// settings.Init(),
		tea.EnterAltScreen,
	)
}

func (m Model) initializeFirstViewCmd() tea.Cmd {
	switch m.currentView {
	case SettingsView:
		return settings.Init()
	case EntriesView:
		return tea.Sequence(
			api.FetchProjects(
				m.config.APIKey,
				m.config.WorkspaceId,
			),
			m.entries.Init(),
		)
	case ReportsView:
		// return m.reports.Init()
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case tea.KeyMsg:
		// Global keybindings
		if !m.showModal {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "1":
				m.currentView = EntriesView
				return m, m.entries.Init()
			case "2":
				m.currentView = ReportsView
				return m, nil
			case "3":
				m.currentView = SettingsView
				return m, nil
			case "?":
				m.showModal = true
				m.modal = modal.NewHelp()
				return m, nil
			}
		}

	case messages.UserLoadedMsg:
		m.userId = msg.UserId
		m.settings, cmd = m.settings.Update(msg)
		return m, cmd

	case messages.ConfigSavedMsg:
		m.config = msg.Config
		m.userId = msg.UserId
		m.workspaceId = msg.WorkspaceId
		_ = m.config.Save()
		return m, nil

	case messages.ProjectsLoadedMsg:
		m.projects = msg.Projects
		m.entries, cmd = m.entries.Update(msg)
		return m, cmd

	case messages.EntrySavedMsg:
		m.showModal = false
		m.entries, cmd = m.entries.Update(msg)
		// return m, cmd
		return m, nil // TEMP

	case messages.EntriesLoadedMsg:
		debug.Log("First Entry %+v", msg.Entries[0])
		m.entries, cmd = m.entries.Update(msg)
		return m, cmd

	case messages.ModalClosedMsg:
		m.showModal = false
		return m, nil
	}

	// Route to modal if showing
	if m.showModal && m.modal != nil {
		*m.modal, cmd = m.modal.Update(msg)
		return m, cmd
	}

	// Route to active view
	switch m.currentView {
	case SettingsView:
		m.settings, cmd = m.settings.Update(msg)
	case EntriesView:
		m.entries, cmd = m.entries.Update(msg)
	case ReportsView:
		// m.reports, cmd = m.reports.Update(msg)
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	var view string

	// Render active view
	switch m.currentView {
	case SettingsView:
		view = m.settings.View()
	case EntriesView:
		view = m.entries.View()
	case ReportsView:
		// view = m.reports.View()
	}

	// Overlay modal if showing
	if m.showModal && m.modal != nil {
		view += modal.Overlay(view, m.modal.View(), m.width, m.height)
	}

	return view
}

func (m Model) Shutdown() tea.Cmd {
	return tea.ExitAltScreen
}
