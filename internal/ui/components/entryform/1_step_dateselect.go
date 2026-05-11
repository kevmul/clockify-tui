package entryform

import (
	"clockify-app/internal/styles"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ================ Date Selection =================
func (m Model) viewDateSelect() string {
	// Implementation of date selection view goes here
	title := styles.TitleStyle.Margin(0, 0).Render("Select Date")
	subtitle := styles.SubtitleStyle.Margin(0, 0, 1, 0).Render("Use arrow keys to navigate, Enter to select")

	// return m.calendar.View()
	dateSelect := fmt.Sprintf("Selected Date\n%s", m.calendar.SelectedDate.Format("Mon, January 02, 2006"))

	return lipgloss.JoinVertical(
		lipgloss.Top,
		title,
		subtitle,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().PaddingRight(2).Render(m.calendar.View().Content),
			dateSelect,
		),
	)

}

func (m Model) updateDateSelect(msg tea.Msg) (Model, tea.Cmd) {
	// Implementation of date selection update goes here
	// Left arrow or 'h' (vim style) - previous day
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {

		case "t":
			// 't' key to jump to today
			m.calendar.SelectedDate = time.Now()
		}
	}
	return m, nil
}
