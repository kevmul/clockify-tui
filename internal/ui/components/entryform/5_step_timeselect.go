package entryform

import (
	"clockify-app/internal/styles"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ================ Time Selection =================
func (m Model) viewTimeInput() string {
	// Implementation of time input view goes here
	title := styles.TitleStyle.Margin(0, 0).Render("Enter Time Range")
	subtitle := styles.SubtitleStyle.Margin(0, 0, 1, 0).Render("Specify the start and end times for your work.")

	return lipgloss.JoinVertical(lipgloss.Top,
		title,
		subtitle,
		fmt.Sprintf("Start Time: %s", m.timeStart.View()),
		fmt.Sprintf("End Time:   %s", m.timeEnd.View()),
		styles.HelpStyle.Render("Press Enter to continue, or Tab/Shift+Tab to navigate."),
	)
}

func (m Model) updateTimeInput(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			// Switch focus between start and end time inputs
			if m.timeStart.Focused() {
				m.timeStart.Blur()
				m.timeEnd.Focus()
			} else {
				m.timeEnd.Blur()
				m.timeStart.Focus()
			}
		case "enter":
			// Move to next step
			m.timeEnd.Blur()
			m.timeStart.Blur()
			m.task.Focus()
			m.step = stepConfirm
		}
	}

	// Update the text inputs
	var cmd tea.Cmd
	m.timeStart, cmd = m.timeStart.Update(msg)
	m.timeEnd, _ = m.timeEnd.Update(msg) // Ignoring second cmd for simplicity

	return m, cmd
}
