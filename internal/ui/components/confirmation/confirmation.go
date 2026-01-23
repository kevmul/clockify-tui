package confirmation

import (
	"clockify-app/internal/messages"
	"clockify-app/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	cursor       int
	itemToDelete string
	itemType     string
}

func New(id string) Model {
	return Model{
		itemToDelete: id,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "left", "right":
			// Toggle between buttons
			if m.cursor == 0 {
				m.cursor = 1
			} else {
				m.cursor = 0
			}

		case "enter":
			// Confirm deletion
			if m.cursor == 0 {
				// Perform deletion logic here
				return m, func() tea.Msg {
					return messages.ItemDeletedMsg{
						ID:   m.itemToDelete,
						Type: m.itemType,
					}
				}
				// You can add actual deletion logic as needed
			} else {
				// Cancel deletion
				return m, func() tea.Msg {
					return messages.ModalClosedMsg{}
				}
			}

		case "esc", "q":
			// Cancel deletion
			return m, nil
		}
	}

	return m, nil
}

func (m Model) View() string {
	title := styles.TitleStyle.Margin(0, 0).Render("Delete Confirmation")
	subtitle := styles.SubtitleStyle.Margin(0, 0, 1, 0).Render("Are you sure you want to delete this entry? This action cannot be undone.")

	return lipgloss.JoinVertical(lipgloss.Top,
		title, subtitle,
		lipgloss.JoinHorizontal(lipgloss.Left,
			m.renderButtons(),
		),
	)
}

func (m Model) renderButtons() string {
	if m.cursor == 0 {
		return lipgloss.JoinHorizontal(lipgloss.Left,
			styles.ActiveButtonStyle.MarginRight(2).Render("Delete"),
			styles.ButtonStyle.Render("Cancel"),
		)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left,
		styles.ButtonStyle.MarginRight(2).Render("Delete"),
		styles.ActiveButtonStyle.Render("Cancel"),
	)

}
