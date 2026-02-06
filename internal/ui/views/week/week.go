package week

import (
	"clockify-app/internal/api"
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"clockify-app/internal/styles"
	"clockify-app/internal/utils"
	"fmt"
	"strings"

	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var TableStyle = lipgloss.NewStyle().Padding(1, 2)

var ColumnWidth = 10

type Model struct {
	config   *config.Config
	entries  []models.Entry
	projects []models.Project

	table  table.Model
	width  int
	height int
	ready  bool
}

func New(cfg *config.Config) Model {
	cols := []table.Column{
		{Title: "Project", Width: 20},
	}
	// Define columns for the week view
	today := time.Now()

	// Find the start of the week (Sunday)
	weekday := int(today.Weekday())
	startOfWeek := today.AddDate(0, 0, -weekday)

	// Now add columns for each day of the week
	for i := range 5 {
		day := startOfWeek.AddDate(0, 0, i+1)
		colTitle := day.Format("Mon01/02 ")
		cols = append(cols, table.Column{Title: colTitle, Width: ColumnWidth})
	}

	// Now we can have the total at the end
	cols = append(cols, table.Column{Title: "Total", Width: ColumnWidth})

	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.Secondary).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(styles.Text).
		Background(styles.Secondary).
		Bold(false)
	t.SetStyles(s)

	return Model{
		config:  cfg,
		entries: []models.Entry{},
		table:   t,
		ready:   false,
	}
}

func (m Model) Init() tea.Cmd {
	today := time.Now()
	startOfWeek := today.AddDate(0, 0, -int(today.Weekday()))
	return api.FetchEntriesForWeek(
		m.config.APIKey,
		m.config.WorkspaceId,
		m.config.UserId,
		startOfWeek,
	)
}

func (m *Model) SetSize(width, height int) {
	h, v := TableStyle.GetFrameSize()
	heightPadding := 6
	widthPadding := 1
	m.width = width
	m.height = height
	if m.ready {
		m.table.SetWidth(m.width - h - widthPadding)
		m.table.SetHeight(m.height - v - heightPadding)
	} else {
		m.table.SetWidth(m.width - h - widthPadding)
		m.table.SetHeight(m.height - v - heightPadding)
	}

	cols := m.table.Columns()
	cols[0].Width = m.width - h - 5 - ColumnWidth*len(cols)
	m.table.SetRows(m.setTableData())
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case messages.EntriesLoadedMsg:
		m.entries = msg.Entries
		m.table.SetRows(m.setTableData())
		m.ready = true

	case messages.ProjectsLoadedMsg:
		m.projects = msg.Projects

	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return TableStyle.Render(m.table.View())
}

func (m Model) setTableData() []table.Row {
	rows := []table.Row{}
	groupedEntries := groupEntriesByProject(m.entries)
	today := time.Now()
	startOfWeek := today.AddDate(0, 0, -int(today.Weekday()))
	dailyTotals := make(map[string]time.Duration)
	for _, group := range groupedEntries {
		project, _ := utils.FindProjectById(m.projects, group[0].ProjectID)

		projectName := project.Name
		if project.ClientName != "" {
			projectName = fmt.Sprintf("%s (%s)", project.Name, project.ClientName)
		}

		row := table.Row{projectName}

		var totalDuration time.Duration = 0
		for i := range 5 {
			day := startOfWeek.AddDate(0, 0, i+1)
			var dayDuration time.Duration = 0

			for _, entry := range group {
				entryDate := entry.TimeInterval.Start
				if entryDate.Year() == day.Year() && entryDate.Month() == day.Month() && entryDate.Day() == day.Day() {
					d := entry.TimeInterval.Duration
					d = strings.TrimPrefix(d, "PT")
					entryDuration, _ := time.ParseDuration(strings.ToLower(d))
					dayDuration += entryDuration
					totalDuration += entryDuration
				}

			}
			dailyTotals[day.Format("2006-01-02")] += dayDuration
			row = append(row, dayDuration.String()) // Placeholder
		}

		dailyTotals["total"] += totalDuration
		row = append(row, totalDuration.String()) // Placeholder
		rows = append(rows, row)
	}
	// Now add the final row for totals
	row := table.Row{"Totals"}
	for i := range 5 {
		day := startOfWeek.AddDate(0, 0, i+1)
		row = append(row, dailyTotals[day.Format("2006-01-02")].String())
	}
	row = append(row, dailyTotals["total"].String())
	rows = append(rows, row)
	return rows
}

// Helper function to group entries by project ID
func groupEntriesByProject(entries []models.Entry) map[string][]models.Entry {
	projectMap := make(map[string][]models.Entry)
	for _, entry := range entries {
		projectMap[entry.ProjectID] = append(projectMap[entry.ProjectID], entry)
	}
	return projectMap
}
