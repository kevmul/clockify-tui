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

	// debug "clockify-app/internal/utils"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	help               *help.Model
	deleteConfirmation *confirmation.Model
	// UI
	viewport viewport.Model
}

func NewEntryForm(cfg *config.Config, projects []models.Project) *Model {
	form := entryform.New(cfg, projects)
	viewport := viewport.New(0, styles.ModalHeight)
	viewport.SetContent(form.View())
	if viewport.Height > viewport.TotalLineCount() {
		viewport.Height = viewport.TotalLineCount()
		viewport.SetContent(form.View())
	}

	return &Model{
		modalType: EntryModal,
		entryForm: &form,
		viewport:  viewport,
	}
}

func UpdateEntryForm(cfg *config.Config, projects []models.Project, entry models.Entry) *Model {
	form := entryform.New(cfg, projects)
	form = form.UpdateEntry(entry)
	viewport := viewport.New(0, styles.ModalHeight)
	viewport.SetContent(form.View())
	return &Model{
		modalType: EntryModal,
		entryForm: &form,
		viewport:  viewport,
	}
}

func NewDeleteConfirmation(entryId string) *Model {
	deleteConfirmation := confirmation.New(entryId, "entry")
	viewport := viewport.New(0, 4)
	viewport.SetContent(deleteConfirmation.View())
	return &Model{
		modalType:          DeleteConfirmation,
		deleteConfirmation: &deleteConfirmation,
		viewport:           viewport,
	}
}

func NewHelp(sections ...help.HelpSection) *Model {
	helpModel := help.New(sections...)
	viewport := viewport.New(0, 10)
	viewport.SetContent(helpModel.View())
	return &Model{
		modalType: HelpModal,
		help:      &helpModel,
		viewport:  viewport,
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
	case tea.KeyMsg:
		// We might move this to the modal themselves later...
		if msg.String() == "esc" || msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, func() tea.Msg {
				return messages.ModalClosedMsg{}
			}
		}
	case messages.TasksLoadedMsg:
		// Pass to entry form if needed
		if m.modalType == EntryModal {
			var cmd tea.Cmd
			*m.entryForm, cmd = m.entryForm.Update(msg)
			m.viewport.SetContent(m.RenderContent())
			return m, cmd
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m.modalType {
	case EntryModal:
		*m.entryForm, cmd = m.entryForm.Update(msg)
		// resize viewportHeight
		if m.entryForm.StepLines <= styles.ModalHeight {
			m.viewport.Height = m.entryForm.StepLines + 1
		} else {
			m.viewport.Height = styles.ModalHeight
		}
	case DeleteConfirmation:
		*m.deleteConfirmation, cmd = m.deleteConfirmation.Update(msg)
	case HelpModal:
		*m.help, cmd = m.help.Update(msg)
	}
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	if _, ok := msg.(tea.KeyMsg); ok {
		// Update viewport content on key events
		m.viewport.SetContent(m.RenderContent())
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {

	viewport := m.viewport.View()

	if m.viewport.TotalLineCount() <= m.viewport.Height {
		// No scrollbar needed
		return styles.ModalStyle.Width(styles.ModalWidth).Render(viewport)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		styles.ModalWithScrollStyle.Width(styles.ModalWidth).Render(viewport),
		utils.RenderScrollbarForModal(m.viewport),
	)

}

func (m Model) RenderContent() string {
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
