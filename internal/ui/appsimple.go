package ui

import (
	"clockify-app/internal/api"
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"clockify-app/internal/styles"
	"clockify-app/internal/ui/components/entryform"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SimpleModel struct {
	// Config and shared state
	config      *config.Config
	userId      string
	workspaceId string
	projects    []models.Project

	// UI
	form   entryform.Model
	width  int
	height int

	// Loading state
	ready bool
}

func NewSimpleModel() SimpleModel {
	cfg, _ := config.LoadConfig()
	return SimpleModel{
		config: cfg,
		form:   entryform.New(cfg, []models.Project{}), // Empty projects for now
	}
}

func (m SimpleModel) Init() tea.Cmd {
	return api.FetchProjects(
		m.config.APIKey,
		m.config.WorkspaceId,
	)
}

func (m SimpleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		// Global key handling can go here if needed
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case messages.ProjectsLoadedMsg:
		m.projects = msg.Projects
		m.form = m.form.SetProjects(m.projects)
		return m, cmd

	case messages.EntrySavedMsg:
		return m, tea.Quit
	}
	m.form, cmd = m.form.Update(msg)
	return m, cmd
}

func (m SimpleModel) View() string {
	form := styles.BoxStyle.Render(m.form.View())
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(form)
}
