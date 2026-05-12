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

	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	scrollOffset int
	title        string
}

func NewEntryForm(cfg *config.Config, projects []models.Project) *Model {
	form := entryform.New(cfg, projects)

	return &Model{
		modalType:    EntryModal,
		entryForm:    &form,
		title:        "New Entry",
		scrollOffset: 0,
	}
}

func UpdateEntryForm(cfg *config.Config, projects []models.Project, entry models.Entry) *Model {
	form := entryform.New(cfg, projects)
	form = form.UpdateEntry(entry)

	return &Model{
		modalType:    EntryModal,
		entryForm:    &form,
		title:        "Edit Entry",
		scrollOffset: 0,
	}
}

func CopyEntryForm(cfg *config.Config, projects []models.Project, entry models.Entry) *Model {
	form := entryform.New(cfg, projects)
	form = form.CopyEntry(entry)

	return &Model{
		modalType:    EntryModal,
		entryForm:    &form,
		title:        "Copy Entry",
		scrollOffset: 0,
	}
}

func NewDeleteConfirmation(entryId string) *Model {
	deleteConfirmation := confirmation.New(entryId, "entry")

	return &Model{
		modalType:          DeleteConfirmation,
		deleteConfirmation: &deleteConfirmation,
		title:              "Confirm Deletion",
		scrollOffset:       0,
	}
}

func NewHelp(sections ...help.HelpSection) *Model {
	helpModel := help.New(sections...)
	return &Model{
		modalType:    HelpModal,
		help:         &helpModel,
		title:        "Help",
		scrollOffset: 0,
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
	case tea.KeyPressMsg:
		switch msg.String() {
		case "j", "down":
			content := m.RenderContent()
			lines := strings.Split(content, "\n")
			maxOffset := max(0, len(lines)-styles.ModalHeight)
			if m.scrollOffset < maxOffset {
				m.scrollOffset++
			}
		case "k", "up":
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}

		case "esc", "q", "ctrl+c":
			var cmd tea.Cmd
			switch m.modalType {
			case EntryModal:
				*m.entryForm, cmd = m.entryForm.Update(msg)
			}
			// Send a message to parent to close the modal
			return m, tea.Batch(cmd, func() tea.Msg {
				return messages.ModalClosedMsg{}
			})
		}
		// We might move this to the modal themselves later...
	case messages.TasksLoadedMsg:
		// Pass to entry form if needed
		if m.modalType == EntryModal {
			var cmd tea.Cmd
			*m.entryForm, cmd = m.entryForm.Update(msg)
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
			// m.viewport.SetHeight(m.entryForm.StepLines + 1)
		} else {
			// m.viewport.SetHeight(styles.ModalHeight)
		}
	case DeleteConfirmation:
		*m.deleteConfirmation, cmd = m.deleteConfirmation.Update(msg)
	case HelpModal:
		*m.help, cmd = m.help.Update(msg)
	}
	cmds = append(cmds, cmd)

	cmds = append(cmds, cmd)

	if _, ok := msg.(tea.KeyPressMsg); ok {
		// Update viewport content on key events
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {

	content := m.RenderContent()
	lines := strings.Split(content, "\n")
	maxOffset := max(0, len(lines)-styles.ModalHeight)
	needsScroll := len(lines) > styles.ModalHeight

	if !needsScroll {
		// No scroll needed - full border
		return tea.NewView(lipgloss.JoinVertical(
			lipgloss.Top,
			createBorderTitle(m.title, styles.ModalWidth, false),
			styles.ModalStyle.Render(content),
		))
	}

	// Scroll needed
	offset := min(m.scrollOffset, maxOffset)
	end := min(offset+styles.ModalHeight, len(lines))
	visibleLines := lines[offset:end]

	if len(visibleLines) < styles.ModalHeight {
		visibleLines = append(visibleLines, "")
	}

	visibleContent := strings.Join(visibleLines, "\n")

	return tea.NewView(lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Top,
			createBorderTitle(m.title, styles.ModalWidth, true),
			styles.ModalWithScrollStyle.Render(visibleContent),
		),
		utils.RenderScrollbarForModal(len(lines), styles.ModalHeight, offset),
	))

	// if m.viewport.TotalLineCount() <= m.viewport.Height() {
	// 	// No scrollbar needed
	// 	return tea.NewView(lipgloss.JoinVertical(
	// 		lipgloss.Top,
	// 		createBorderTitle(m.title, styles.ModalWidth, false),
	// 		styles.ModalStyle.Render(content),
	// 	))
	// }
	//
	// viewport := lipgloss.JoinVertical(
	// 	lipgloss.Top,
	// 	createBorderTitle(m.title, styles.ModalWidth, true),
	// 	styles.ModalWithScrollStyle.Render(content),
	// )
	//
	// return tea.NewView(lipgloss.JoinHorizontal(
	// 	lipgloss.Top,
	// 	viewport,
	// 	utils.RenderScrollbarForModal(m.viewport),
	// ))

}

func createBorderTitle(title string, modalWidth int, withScroll bool) string {
	borderChar := styles.CustomBorder.Top
	titleLength := lipgloss.Width(title)
	if titleLength >= modalWidth-2 {
		// Title is too long to fit, return it as is (it will be truncated by the modal)
		return title
	}

	leftBorderLength := 2                                                //
	rightBorderLength := modalWidth - titleLength - leftBorderLength - 4 // 2 for the spaces around the title

	s := styles.CustomBorder.TopLeft +
		strings.Repeat(string(borderChar), leftBorderLength) +
		" " + title + " " +
		strings.Repeat(string(borderChar), rightBorderLength)

	if !withScroll {
		s += styles.CustomBorder.TopRight
	} else {
		s += borderChar
	}

	return styles.ModalTitleStyle.Render(s)

}

func (m Model) RenderContent() string {
	switch m.modalType {
	case EntryModal:
		return m.entryForm.View().Content
	case DeleteConfirmation:
		return m.deleteConfirmation.View().Content
	case HelpModal:
		return m.help.View().Content
	}
	return "MODAL"
}
