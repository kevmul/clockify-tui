package project

import (
	"clockify-app/internal/api"
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"clockify-app/internal/styles"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// A single page view

type Model struct {
	config  *config.Config
	project models.Project
	tasks   []models.Task
	ready   bool
}

func New(cfg *config.Config, project models.Project, tasks []models.Task) Model {
	return Model{
		config:  cfg,
		project: project,
		tasks:   tasks,
		ready:   false,
	}
}

func (m Model) Init() tea.Cmd {
	return api.FetchTasks(m.config.APIKey, m.config.WorkspaceId, m.project.ID)
}

func (m Model) Update(msg any) (Model, tea.Cmd) {
	// var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n": // Create new task
			// cmd = api.CreateTask(m.config.APIKey, m.config)
			// cmds = append(cmds, cmd)

		case "esc", "b": // Go back to projects view
			m.ready = false
			cmds = append(cmds, func() tea.Msg { return messages.ExitViewMsg{} })
		}

	case messages.TasksLoadedMsg:
		m.tasks = msg.Tasks
		m.ready = true
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	title := styles.TitleStyle.Render("Project: " + m.project.Name)
	if !m.ready {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			title,
			"\nLoading tasks...",
		)
	}
	if len(m.tasks) == 0 {
		return lipgloss.JoinVertical(
			lipgloss.Top,
			title,
			"\nNo tasks found for this project.",
		)
	}
	s := strings.Builder{}

	for _, task := range m.tasks {
		s.WriteString(fmt.Sprintf("- %s\n", task.Name))

	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		title,
		s.String(),
	)
}
