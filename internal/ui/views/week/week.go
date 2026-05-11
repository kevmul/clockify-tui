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

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

var TableStyle = lipgloss.NewStyle().Padding(1, 2)

var ColumnWidth = 11

type Model struct {
	config          *config.Config
	entries         []models.Entry
	projects        []models.Project
	table           *table.Table
	weekStart       time.Time
	projectColWidth int
	width           int
	height          int
	ready           bool
}

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(styles.Primary).
			Bold(true).
			Align(lipgloss.Center)

	cellStyle = lipgloss.NewStyle().Padding(0, 1).Align(lipgloss.Right)

	totalColStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Width(ColumnWidth).
			Foreground(styles.Secondary).
			Align(lipgloss.Right)
)

func New(cfg *config.Config) Model {
	today := time.Now()
	weekday := int(today.Weekday())
	startOfWeek := today.AddDate(0, 0, -weekday)

	m := Model{
		config:    cfg,
		entries:   []models.Entry{},
		weekStart: startOfWeek,
		ready:     false,
	}

	m.table = table.New().
		BorderStyle(lipgloss.NewStyle().Foreground(styles.Secondary)).
		StyleFunc(func(row, col int) lipgloss.Style {
			numCols := 7 // Project + 5 days + Total
			if row == table.HeaderRow {
				if col == 0 {
					return headerStyle.Width(m.projectColWidth)
				}
				if col == numCols-1 {
					// Last column is always the Total col
					return headerStyle.Foreground(styles.Secondary)
				}
				return headerStyle
			}
			// Last column is the Totals column
			if col == numCols-1 {
				return totalColStyle
			}
			style := cellStyle
			if row%2 == 0 {
				return style.Foreground(styles.Muted)
			}
			if col == 0 {
				style.Width(m.projectColWidth).Align(lipgloss.Left)
			}

			return style
		})

	return m
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
	m.width = width
	m.height = height
	frameWidth, _ := TableStyle.GetFrameSize()
	m.projectColWidth = width - frameWidth - (ColumnWidth * 6) - 2
}

func (m *Model) PreviousWeek() tea.Cmd {
	m.weekStart = m.weekStart.AddDate(0, 0, -7)
	m.ready = false
	return api.FetchEntriesForWeek(
		m.config.APIKey,
		m.config.WorkspaceId,
		m.config.UserId,
		m.weekStart,
	)
}

func (m *Model) NextWeek() tea.Cmd {
	m.weekStart = m.weekStart.AddDate(0, 0, 7)
	m.ready = false
	return api.FetchEntriesForWeek(
		m.config.APIKey,
		m.config.WorkspaceId,
		m.config.UserId,
		m.weekStart,
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "h", "left":
			cmds = append(cmds, m.PreviousWeek())
		case "l", "right":
			cmds = append(cmds, m.NextWeek())
		}

	case messages.EntriesLoadedMsg:
		m.entries = msg.Entries
		m.table.ClearRows()
		m.table.Headers(m.tableHeaders()...)
		m.table.Rows(m.setTableData()...)
		m.ready = true

	case messages.ProjectsLoadedMsg:
		m.projects = msg.Projects
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	return tea.NewView(TableStyle.Render(m.table.Render()))
}

func (m Model) tableHeaders() []string {
	headers := []string{"Project"}
	for i := range 5 {
		day := m.weekStart.AddDate(0, 0, i+1)
		headers = append(headers, day.Format("Mon 01/02"))
	}
	headers = append(headers, "Total")
	return headers
}

func (m Model) setTableData() [][]string {
	rows := [][]string{}
	groupedEntries := groupEntriesByProject(m.entries)
	startOfWeek := m.weekStart
	dailyTotals := make(map[string]time.Duration)

	for _, group := range groupedEntries {
		project, _ := utils.FindProjectById(m.projects, group[0].ProjectID)

		projectName := project.Name
		if project.ClientName != "" {
			projectName = fmt.Sprintf("%s (%s)", project.Name, project.ClientName)
		}

		row := []string{projectName}
		var totalDuration time.Duration

		for i := range 5 {
			day := startOfWeek.AddDate(0, 0, i+1)
			var dayDuration time.Duration

			for _, entry := range group {
				entryDate := entry.TimeInterval.Start
				if entryDate.Year() == day.Year() &&
					entryDate.Month() == day.Month() &&
					entryDate.Day() == day.Day() {
					d := strings.TrimPrefix(entry.TimeInterval.Duration, "PT")
					entryDuration, _ := time.ParseDuration(strings.ToLower(d))
					dayDuration += entryDuration
					totalDuration += entryDuration
				}
			}

			dailyTotals[day.Format("2006-01-02")] += dayDuration
			row = append(row, formatDuration(dayDuration))
		}

		dailyTotals["total"] += totalDuration
		row = append(row, formatDuration(totalDuration))
		rows = append(rows, row)
	}

	// Totals row
	totalsRow := []string{"Totals"}
	for i := range 5 {
		day := startOfWeek.AddDate(0, 0, i+1)
		totalsRow = append(totalsRow, formatDuration(dailyTotals[day.Format("2006-01-02")]))
	}
	totalsRow = append(totalsRow, formatDuration(dailyTotals["total"]))
	rows = append(rows, totalsRow)

	return rows
}

func groupEntriesByProject(entries []models.Entry) map[string][]models.Entry {
	projectMap := make(map[string][]models.Entry)
	for _, entry := range entries {
		projectMap[entry.ProjectID] = append(projectMap[entry.ProjectID], entry)
	}
	return projectMap
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "-"
	}
	h := int(d.Hours())
	min := int(d.Minutes()) % 60
	if h == 0 {
		return fmt.Sprintf("%dm", min)
	}
	return fmt.Sprintf("%dh %dm", h, min)
}
