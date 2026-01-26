package modal

import (
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"clockify-app/internal/styles"
	"clockify-app/internal/ui/components/confirmation"
	"clockify-app/internal/ui/components/entryform"
	"clockify-app/internal/ui/components/help"
	"clockify-app/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type ModalType int

const (
	EntryModal ModalType = iota
	DeleteConfirmation
	HelpModal
)

type Model struct {
	modalType          ModalType
	entryForm          *entryform.Model
	deleteConfirmation *confirmation.Model
	// UI
	width, height int
	help          *help.Model
}

func NewEntryForm(cfg *config.Config, projects []models.Project) *Model {
	form := entryform.New(cfg, projects)
	return &Model{
		modalType: EntryModal,
		entryForm: &form,
	}
}

func UpdateEntryForm(cfg *config.Config, projects []models.Project, entry models.Entry) *Model {
	form := entryform.New(cfg, projects)
	form = form.UpdateEntry(entry)
	return &Model{
		modalType: EntryModal,
		entryForm: &form,
	}
}

func NewDeleteConfirmation(entryId string) *Model {
	deleteConfirmation := confirmation.New(entryId, "entry")
	return &Model{
		modalType:          DeleteConfirmation,
		deleteConfirmation: &deleteConfirmation,
	}
}

func NewHelp(sections ...help.HelpSection) *Model {

	helpModel := help.New(sections...)

	return &Model{
		modalType: HelpModal,
		help:      &helpModel,
	}
}

func (m Model) Init() tea.Cmd {
	switch m.modalType {
	case EntryModal:
		return m.entryForm.Init()
	case DeleteConfirmation:
		return m.deleteConfirmation.Init()
	case HelpModal:
		return m.help.Init()
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		// We might move this to the modal themselves later...
		if msg.String() == "esc" || msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, func() tea.Msg {
				return messages.ModalClosedMsg{}
			}
		}
	}

	var cmd tea.Cmd
	switch m.modalType {
	case EntryModal:
		*m.entryForm, cmd = m.entryForm.Update(msg)
	case DeleteConfirmation:
		*m.deleteConfirmation, cmd = m.deleteConfirmation.Update(msg)
	case HelpModal:
		*m.help, cmd = m.help.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	switch m.modalType {
	case EntryModal:
		return m.entryForm.View()
	case DeleteConfirmation:
		return m.deleteConfirmation.View()
	case HelpModal:
		return m.help.View()
	}
	return "MODAL"
}

// Overlay renders a modal on top of existing content
func Overlay(base, modal string, width, height int) string {
	return utils.RenderWithModal(height, width, base, styles.ModalStyle.Render(modal))
}
