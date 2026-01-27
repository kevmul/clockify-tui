package entryform

import (
	"clockify-app/internal/styles"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ================ Date Selection =================
func (m Model) viewDateSelect() string {
	// Implementation of date selection view goes here
	title := styles.TitleStyle.Margin(0, 0).Render("Select Date")
	subtitle := styles.SubtitleStyle.Margin(0, 0, 1, 0).Render("Use arrow keys to navigate, Enter to select")

	dateSelect := fmt.Sprintf("📅 %s", m.date.Format("Monday, January 2, 2006"))
	if m.date.Equal(time.Now()) {
		dateSelect += " (Today)"
	}
	return lipgloss.JoinVertical(lipgloss.Top, title, subtitle, dateSelect)
}

func (m Model) updateDateSelect(msg tea.Msg) (Model, tea.Cmd) {
	// Implementation of date selection update goes here
	// Left arrow or 'h' (vim style) - previous day
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "t":
			// 't' key to jump to today
			m.date = time.Now()

		case "left", "h":
			if m.step == stepDateSelect {
				m.date = m.date.AddDate(0, 0, -1)
			}

		// Right arrow or 'l' (vim style) - next day
		case "right", "l":
			if m.step == stepDateSelect {
				m.date = m.date.AddDate(0, 0, 1)
			}

		case "enter":
			// Move to next step
			m.description.Focus()
			m.step = stepDescriptionInput
		}
	}
	return m, nil
}
