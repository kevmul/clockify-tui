package entries

import (
	"clockify-app/internal/api"
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"clockify-app/internal/utils"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	config   *config.Config
	projects []models.Project
	tasks    []models.Task
	entries  []models.Entry
	cursor   int

	list list.Model
}

func New(cfg *config.Config) Model {
	list := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Clockify Entries"
	list.SetShowStatusBar(true)
	list.SetFilteringEnabled(true)
	list.SetShowHelp(false)

	return Model{
		config:  cfg,
		entries: []models.Entry{},
		list:    list,
	}
}

func (m Model) Init() tea.Cmd {
	return api.FetchEntries(
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
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-4)

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			// Handle down navigation
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}
		case "k", "up":
			// Handle up navigation
			if m.cursor > 0 {
				m.cursor--
			}
		}

	case messages.EntriesLoadedMsg:
		m.entries = msg.Entries
		items := make([]list.Item, len(m.entries))
		for i, entry := range m.entries {
			// Get the description or a placeholder
			description := entry.Description
			if description == "" {
				description = "(No Description)"
			}
			// Get the project name or a default
			projectName := "No Project"
			project, _ := utils.FindProjectById(m.projects, entry.ProjectID)
			if project.ID != "" {
				projectName = fmt.Sprintf("%s - %s", project.Name, project.ClientName)
			}
			items[i] = item{
				title: description,
				desc: fmt.Sprintf(
					"%s  %s - %s (%s)",

					entry.TimeInterval.Start.In(time.Local).Format("Jan 02 2006"),
					entry.TimeInterval.Start.In(time.Local).Format("3:04PM"),
					entry.TimeInterval.End.In(time.Local).Format("3:04PM"),
					projectName,
				),
			}
		}
		m.list.SetItems(items)

	case messages.ProjectsLoadedMsg:
		m.projects = msg.Projects

	case messages.ErrorMsg:
		// Handle error (could set an error field in the model)
		fmt.Printf("Error: %v\n", msg.Err)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

var docStyle = lipgloss.NewStyle()

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m Model) View() string {

	return docStyle.Render(m.list.View())

}
