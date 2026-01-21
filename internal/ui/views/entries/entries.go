package entries

import (
	"clockify-app/internal/api"
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"clockify-app/internal/utils"
	"fmt"

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
	return Model{
		config:  cfg,
		entries: []models.Entry{},
		list:    list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
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
		m.list.SetSize(msg.Width-h, msg.Height-v)

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
			projectName := "No Project"
			project, _ := utils.FindProjectById(m.projects, entry.ProjectID)
			if project.ID != "" {
				projectName = fmt.Sprintf("%s - %s", project.Name, project.ClientName)
			}
			items[i] = item{
				title: entry.Description,
				desc: fmt.Sprintf(
					"%s  %s - %s (%s)",
					entry.TimeInterval.Start.Format("Jan 02 2006"),
					entry.TimeInterval.Start.Format("3:04PM"),
					entry.TimeInterval.End.Format("3:04PM"),
					projectName,
				),
			}
		}
		m.list.SetItems(items)
		m.list.Title = "Time Entries"

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

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m Model) View() string {

	return docStyle.Render(m.list.View())

}
