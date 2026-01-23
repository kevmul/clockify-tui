package reports

import (
	"clockify-app/internal/api"
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	debug "clockify-app/internal/utils"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	cursor int
	config *config.Config

	list list.Model
}

func New(cfg *config.Config) Model {
	list := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Clockify Entries"
	list.SetShowStatusBar(true)
	list.SetFilteringEnabled(true)
	list.SetShowHelp(false)

	return Model{
		config: cfg,
		list:   list,
	}
}

func (m Model) Init() tea.Cmd {
	debug.Log("Initializing Reports Model")
	return api.FetchReports(
		m.config.APIKey,
		m.config.WorkspaceId,
		m.config.UserId,
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	// Set the list size when the window size changes
	case messages.ReportsLoadedMsg:
		debug.Log("Reports loaded:", len(msg.Reports))
		items := make([]list.Item, len(msg.Reports))
		for i, report := range msg.Reports {
			items[i] = item{
				title: report.Name,
				desc:  report.Type,
			}
		}
		m.list.SetItems(items)
	}

	// Let the list component handle its own updates
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m Model) View() string {

	return lipgloss.NewStyle().Render(m.list.View())
}
