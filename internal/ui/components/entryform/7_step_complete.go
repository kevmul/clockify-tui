package entryform

import (
	"clockify-app/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewCompletionInput() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		styles.SuccessStyle.Render("Time entry created successfully!"),
		styles.SubtitleStyle.Render("Press [enter] to close."),
	)
}

func (m Model) updateComplete(msg tea.Msg) (Model, tea.Cmd) {
	// Implementation of description input update goes here
	return m, nil
}
