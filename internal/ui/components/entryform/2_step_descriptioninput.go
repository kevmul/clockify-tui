package entryform

import (
	"clockify-app/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewDescriptionInput() string {
	title := styles.TitleStyle.MarginBottom(0).Render("Enter Description")
	subtitle := styles.SubtitleStyle.MarginBottom(1).Render("Provide a brief description of the work done.")
	return lipgloss.JoinVertical(
		lipgloss.Top,
		title,
		subtitle,
		m.description.View(),
		styles.HelpStyle.Render("Press Enter to continue, or Tab/Shift+Tab to navigate."),
	)
}

func (m Model) updateDescriptionInput(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Move to next step or finish
			m.description.Blur()
			m.step++
		}
	}
	m.description, cmd = m.description.Update(msg)

	return m, cmd
}
