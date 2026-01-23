package entryform

import (
	"clockify-app/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ================ Task Selection =================
func (m Model) viewTaskInput() string {
	// Implementation of task input view goes here
	title := styles.TitleStyle.Margin(0, 0).Render("Enter Task Description")
	subtitle := styles.SubtitleStyle.Margin(0, 0, 1, 0).Render("Describe the work you did during this time period.")

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		subtitle,
		m.taskName.View(),
		styles.HelpStyle.Render("Press Enter to continue, or Tab/Shift+Tab to navigate."),
	)
}

func (m Model) updateTaskInput(msg tea.Msg) (Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Move to next step or finish
			m.taskName.Blur()
			m.step++
		}
	}
	m.taskName, cmd = m.taskName.Update(msg)

	return m, cmd
}
