package modal

import (
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"clockify-app/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ModalType int

const (
	EntryModal ModalType = iota
	HelpModal
)

type Model struct {
	modalType ModalType
	// entryForm *entryform.Model
	// help      *help.Model
}

func NewEntryFrom(projects []models.Project, tasks []models.Task) *Model {
	// form := entryform.New(projects, tasks)
	return &Model{
		modalType: EntryModal,
		// entryForm: form,
	}
}

func NewHelp() *Model {
	// helpModel := help.New()
	return &Model{
		modalType: HelpModal,
		// help:      helpModel,
	}
}

func (m Model) Init() tea.Cmd {
	switch m.modalType {
	case EntryModal:
		// return m.entryForm.Init()
	case HelpModal:
		// return m.help.Init()
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			return m, func() tea.Msg {
				return messages.ModalClosedMsg{}
			}
		}
	}

	var cmd tea.Cmd
	switch m.modalType {
	case EntryModal:
		// *m.entryForm, cmd = m.entryForm.Update(msg)
	case HelpModal:
		// *m.help, cmd = m.help.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	switch m.modalType {
	case EntryModal:
		// return ui.ModalStyle.Render(m.entryForm.View())
	case HelpModal:
		// return ui.ModalStyle.Render(m.help.View())
	}
	return "MODAL"
}

// Overlay renders a modal on top of existing content
func Overlay(base, modal string, width, height int) string {
	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 0).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)
	// return lipgloss.Place(
	// 	width, height,
	// 	lipgloss.Center,
	// 	lipgloss.Center,
	// 	base+dialogBoxStyle.Render(modal),
	// 	lipgloss.WithWhitespaceChars("░"),
	// 	// lipgloss.WithWhitespaceChars(" "),
	// 	lipgloss.WithWhitespaceForeground(lipgloss.Color("#1a1a1a")),
	// )
	return utils.PlaceOverlay(width, height, dialogBoxStyle.Render(modal), base, true)
}
