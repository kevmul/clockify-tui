package month

import (
	"clockify-app/internal/api"
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"clockify-app/internal/styles"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

var TableStyle = lipgloss.NewStyle().Padding(0, 2)

var ColumnWidth = 10

type Model struct {
	config       *config.Config
	entries      []models.Entry
	currentMonth time.Time

	table *table.Table

	width  int
	height int
	ready  bool
}

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(styles.Primary).
			Bold(true).
			Align(lipgloss.Center)

	cellStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Width(14).
			Align(lipgloss.Center)
)

func New(cfg *config.Config) Model {
	m := Model{
		config:       cfg,
		entries:      []models.Entry{},
		currentMonth: time.Now(),
		ready:        false,
	}

	m.table = table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(styles.Secondary)).
		BorderRow(true).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			return cellStyle
		})

	return m
}

func (m Model) Init() tea.Cmd {
	return api.FetchEntriesForMonth(m.config.APIKey, m.config.WorkspaceId, m.config.UserId, m.currentMonth)
}

func (m Model) View() tea.View {

	footer := m.renderFooter()

	return tea.NewView(lipgloss.JoinVertical(
		lipgloss.Left,
		styles.TitleStyle.
			PaddingTop(1).
			PaddingLeft(2).
			Render(m.currentMonth.Format("January")),
		TableStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.table.Render(),
				footer,
			),
		),
	))
}

func (m Model) renderFooter() string {
	monthTotal := m.calculateMonthTotal()

	totalStyle := lipgloss.NewStyle()

	label := lipgloss.NewStyle().
		Foreground(styles.Secondary).
		Padding(0, 1).
		Bold(true).
		Render("Month Total")

	value := lipgloss.NewStyle().
		Foreground(styles.Text).
		Render(formatDuration(monthTotal))

	content := lipgloss.JoinHorizontal(
		lipgloss.Left,
		label,
		"  ",
		value,
	)

	return totalStyle.Render(content)
}

func (m Model) calculateMonthTotal() time.Duration {
	var total time.Duration
	for _, entry := range m.entries {
		if entry.TimeInterval.End.IsZero() {
			continue
		}
		total += entry.TimeInterval.End.Sub(entry.TimeInterval.Start)
	}
	return total
}

func (m Model) NextMonth() (Model, tea.Cmd) {
	m.currentMonth = m.currentMonth.AddDate(0, 1, 0)
	m.ready = false
	m.entries = []models.Entry{}
	m.table.ClearRows()
	return m, api.FetchEntriesForMonth(m.config.APIKey, m.config.WorkspaceId, m.config.UserId, m.currentMonth)
}

func (m Model) PreviousMonth() (Model, tea.Cmd) {
	m.currentMonth = m.currentMonth.AddDate(0, -1, 0)
	m.ready = false
	m.entries = []models.Entry{}
	m.table.ClearRows()
	return m, api.FetchEntriesForMonth(m.config.APIKey, m.config.WorkspaceId, m.config.UserId, m.currentMonth)

}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "h", "left":
			m, cmd = m.PreviousMonth()
			cmds = append(cmds, cmd)
		case "l", "right":
			m, cmd = m.NextMonth()
			cmds = append(cmds, cmd)
		}

	case messages.EntriesLoadedMsg:
		m.entries = msg.Entries
		m.table.ClearRows()
		m.table.Headers(m.tableHeaders()...)
		m.table.Rows(m.setTableData()...)
		m.SetSize(m.width, m.height)
		m.ready = true
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

var WeekColumnWidth = 8
var WeekLabelWidth = 8

func (m Model) tableHeaders() []string {
	return []string{"Mon", "Tues", "Wed", "Thurs", "Fri", "Total"}
}

func (m Model) setTableData() [][]string {
	// Aggregate daily totals from entries
	dailyTotals := make(map[string]time.Duration)
	for _, entry := range m.entries {
		if entry.TimeInterval.End.IsZero() {
			continue
		}
		duration := entry.TimeInterval.End.Sub(entry.TimeInterval.Start)
		day := entry.TimeInterval.Start.Format("2006-01-02")
		dailyTotals[day] += duration
	}

	startOfMonth := time.Date(m.currentMonth.Year(), m.currentMonth.Month(), 1, 0, 0, 0, 0, time.Local)
	daysInMonth := time.Date(m.currentMonth.Year(), m.currentMonth.Month()+1, 0, 0, 0, 0, 0, time.Local).Day()

	// Group days into Mon–Fri weeks
	type week struct {
		label string
		days  [5]time.Time
	}

	var weeks []week
	var current *week
	weekNum := 0

	for i := range daysInMonth {
		day := startOfMonth.AddDate(0, 0, i)
		weekday := day.Weekday()

		if weekday == time.Saturday || weekday == time.Sunday {
			continue
		}

		if weekday == time.Monday || current == nil {
			weekNum++
			weeks = append(weeks, week{})
			current = &weeks[len(weeks)-1]
		}

		idx := int(weekday) - 1 // Mon=0, Tue=1, ..., Fri=4
		current.days[idx] = day
	}

	// Build rows
	rows := [][]string{}
	for _, w := range weeks {
		var weekTotal time.Duration
		row := []string{}
		for _, day := range w.days {
			if day.IsZero() {
				row = append(row, "\n")
				continue
			}
			key := day.Format("2006-01-02")
			d := dailyTotals[key]
			weekTotal += d
			date := day.Format("01/02")
			row = append(row, fmt.Sprintf("%s\n%s", date, formatDuration(d)))
		}
		row = append(row, fmt.Sprintf("\n%s", formatDuration(weekTotal)))
		rows = append(rows, row)
	}

	return rows
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "-"
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h == 0 {
		return fmt.Sprintf("%dm", m)
	}
	if m == 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dh %dm", h, m)
}
