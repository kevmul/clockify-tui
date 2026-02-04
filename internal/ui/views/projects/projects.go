package projects

import (
	"clockify-app/internal/api"
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"clockify-app/internal/styles"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	config   *config.Config
	projects []models.Project
	list     list.Model
	ready    bool
	width    int
	height   int
}

var docStyle = lipgloss.NewStyle()

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func New(cfg *config.Config) Model {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(styles.Primary).BorderLeftForeground(styles.Primary)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(styles.Tertiary).BorderLeftForeground(styles.Primary)

	list := list.New([]list.Item{}, d, 0, 0)

	list.SetShowTitle(false)
	list.SetShowStatusBar(true)
	list.SetFilteringEnabled(true)
	list.SetShowHelp(false)

	return Model{
		config:   cfg,
		projects: []models.Project{},
		list:     list,
		ready:    false,
	}
}

func (m Model) Init() tea.Cmd {
	return api.FetchProjects(
		m.config.APIKey,
		m.config.WorkspaceId,
	)
}

func (m Model) Update(msg any) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	// Set the list size when the window size changes
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-5)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter": // Open selected project
			if len(m.projects) > 0 {
				selectedProject := m.projects[m.list.Index()]
				return m, func() tea.Msg {
					return messages.ProjectSelectedMsg{
						Project: selectedProject,
					}
				}
			}
		}

	case messages.ProjectsLoadedMsg:
		m.projects = msg.Projects
		items := make([]list.Item, len(m.projects))
		for i, project := range m.projects {
			title := project.Name
			if project.ClientName != "" {
				title = fmt.Sprintf("%s (%s)", project.Name, project.ClientName)
			}

			items[i] = item{
				title: title,
				desc:  project.ID,
			}
		}
		m.list.SetItems(items)
		m.ready = true
		return m, nil
	}
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "Loading projects..."
	}

	if len(m.projects) == 0 {
		return "No projects found."
	}
	return docStyle.Render(m.list.View())
}
