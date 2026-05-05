package month

import (
	"clockify-app/internal/api"
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"clockify-app/internal/styles"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var TableStyle = lipgloss.NewStyle().Padding(0, 2)

var ColumnWidth = 10

type Model struct {
	config       *config.Config
	entries      []models.Entry
	currentMonth time.Time

	table        table.Model

	width        int 
	height       int 
	ready        bool
}

func New(cfg *config.Config) Model {
	m := Model {
		config: cfg,
		entries: []models.Entry{},
		currentMonth: time.Now(),
		ready: false,
	}

	t := table.New(
		table.WithColumns(m.setTableColumns()),
		table.WithFocused(false),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.Secondary).
		BorderBottom(true).
		Bold(false)
	s.Selected = lipgloss.NewStyle()
	t.SetStyles(s)

	m.table = t

	return m
}

func (m Model) Init() tea.Cmd {
	return api.FetchEntriesForMonth(m.config.APIKey, m.config.WorkspaceId, m.config.UserId, m.currentMonth)
}

func (m Model) View() string {

	footer := m.renderFooter()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		styles.TitleStyle.
			PaddingTop(1).
			PaddingLeft(2).
			Render(m.currentMonth.Format("January")),
		TableStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.table.View(),
				footer,
			),
		),
	)
}

func (m Model) renderFooter() string {
	monthTotal := m.calculateMonthTotal()

	totalStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.Secondary).
		BorderTop(true).
		Width(m.table.Width())

	label := lipgloss.NewStyle().
		Foreground(styles.Secondary).
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
	m.table.SetRows([]table.Row{})
	return m, api.FetchEntriesForMonth(m.config.APIKey, m.config.WorkspaceId, m.config.UserId, m.currentMonth)
}

func (m Model) PreviousMonth() (Model, tea.Cmd) {
	m.currentMonth = m.currentMonth.AddDate(0, -1, 0)
	m.ready = false
	m.entries = []models.Entry{}
	m.table.SetRows([]table.Row{})
	return m, api.FetchEntriesForMonth(m.config.APIKey, m.config.WorkspaceId, m.config.UserId, m.currentMonth)

}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd 
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
		m.table.SetColumns(m.setTableColumns())
		m.table.SetRows(m.setTableData())
		m.SetSize(m.width, m.height)
		m.ready = true
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) SetSize(width, height int) {
	h, _ := TableStyle.GetFrameSize()
	widthPadding := 1
	m.width = width
	m.height = height
	if m.ready {
		m.table.SetWidth(m.width - h - widthPadding)
		m.table.SetHeight(7)
	} else {
		m.table.SetWidth(m.width - h - widthPadding)
		m.table.SetHeight(7)
	}

	//cols := m.table.Columns()
	//cols[0].Width = m.width - h - 5 - ColumnWidth*len(cols)
	//m.table.SetRows(m.setTableData())
}

var WeekColumnWidth = 8
var WeekLabelWidth = 8

func (m Model) setTableColumns() []table.Column {
	cols := []table.Column{
		{Title: "", Width: WeekLabelWidth},
		{Title: "Mon", Width: WeekColumnWidth},
		{Title: "Tue", Width: WeekColumnWidth},
		{Title: "Wed", Width: WeekColumnWidth},
		{Title: "Thu", Width: WeekColumnWidth},
		{Title: "Fri", Width: WeekColumnWidth},
		{Title: "Total", Width: WeekColumnWidth},
	}
	return cols
}

func (m Model) setTableData() []table.Row {
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
		days  [5]time.Duration // Mon=0 ... Fri=4
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
			weeks = append(weeks, week{label: fmt.Sprintf("Week %d", weekNum)})
			current = &weeks[len(weeks)-1]
		}

		idx := int(weekday) - 1 // Mon=0, Tue=1, ..., Fri=4
		key := day.Format("2006-01-02")
		current.days[idx] = dailyTotals[key]
	}

	// Build rows
	rows := []table.Row{}
	for _, w := range weeks {
		var weekTotal time.Duration
		row := table.Row{w.label}
		for _, d := range w.days {
			weekTotal += d
			row = append(row, formatDuration(d))
		}
		row = append(row, formatDuration(weekTotal))
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
