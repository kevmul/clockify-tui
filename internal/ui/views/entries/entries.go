package entries

import (
	"clockify-app/internal/api"
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"clockify-app/internal/styles"
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
	entries  []models.Entry

	list list.Model
}

func New(cfg *config.Config) Model {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(styles.Primary).BorderLeftForeground(styles.Primary)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(styles.Tertiary).BorderLeftForeground(styles.Primary)

	list := list.New([]list.Item{}, d, 0, 0)
	list.Title = "Clockify Entries"
	list.SetShowTitle(false)
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
		m.list.SetSize(msg.Width-h, msg.Height-v-5)

	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			// Edit the selected entry
			if len(m.entries) > 0 {
				selectedEntry := m.entries[m.list.Index()]
				// Open the edit modal (not implemented here)
				return m, func() tea.Msg {
					return messages.EntryUpdateStartedMsg{Entry: selectedEntry}
				}
			}
		case "d":
			// Delete the selected entry
			if len(m.entries) > 0 {
				selectedEntry := m.entries[m.list.Index()]
				// Open the delete confirmation modal (not implemented here)
				return m, func() tea.Msg {
					return messages.EntryDeleteStartedMsg{EntryId: selectedEntry.ID}
				}
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
			if project.ID != "" && project.ClientName != "" {
				projectName = fmt.Sprintf("%s - %s", project.Name, project.ClientName)
			}
			if project.ID != "" {
				projectName = fmt.Sprintf("%s", project.Name)
			}
			items[i] = item{
				title: description,
				desc: fmt.Sprintf(
					"%s-%s (%s)",

					// entry.TimeInterval.Start.In(time.Local).Format("Mon, Jan 02 2006 3:04PM"),
					entry.TimeInterval.Start.In(time.Local).Format("Mon, 2006_01_02 03:04PM"),
					entry.TimeInterval.End.In(time.Local).Format("03:04PM"),
					projectName,
				),
			}
		}
		m.list.SetItems(items)

	case messages.ProjectsLoadedMsg:
		m.projects = msg.Projects

	case messages.ItemDeletedMsg:
		if msg.Type == "entry" {
			c := api.NewClient(m.config.APIKey)
			err := c.DeleteTimeEntry(m.config.WorkspaceId, msg.ID)
			if err != nil {
				return m, func() tea.Msg {
					return messages.ErrorMsg{Err: err}
				}
			}
			return m, api.FetchEntries(
				m.config.APIKey,
				m.config.WorkspaceId,
				m.config.UserId,
			)
		}
		return m, nil

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
