package entryform

import (
	"clockify-app/internal/styles"
	"clockify-app/internal/utils"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	timeStartErr = ""
	timeEndErr   = ""
)

// ================ Time Selection =================
func (m Model) viewTimeInput() string {
	// Implementation of time input view goes here
	title := styles.TitleStyle.Margin(0, 0).Render("Enter Time Range")
	subtitle := styles.SubtitleStyle.Margin(0, 0, 1, 0).Render("Specify the start and end times for your work.")

	lines := []string{
		title,
		subtitle,
		fmt.Sprintf("Start Time: %s", m.timeStart.View()),
	}

	if timeStartErr != "" {
		lines = append(lines, styles.ErrorStyle.Render(timeStartErr))
	}

	lines = append(lines, fmt.Sprintf("End Time:   %s", m.timeEnd.View()))

	if timeEndErr != "" {
		lines = append(lines, styles.ErrorStyle.Render(timeEndErr))
	}

	lines = append(lines, styles.HelpStyle.Render("Press Enter to continue, or Tab/Shift+Tab to navigate."))

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lines...,
	)
}

func (m Model) updateTimeInput(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Switch focus between start and end time inputs
			if m.timeStart.Focused() {
				m.timeStart.Blur()
				m.timeEnd.Focus()
			} else {
				m.timeEnd.Blur()
				m.timeStart.Focus()
			}
		case "enter":
			// Move to next step
			timeStartErr = ""
			timeEndErr = ""
			m.validate()
			if timeStartErr != "" || timeEndErr != "" {
				// Show errors
				return m, nil
			}

			// Move to next step
			m.timeEnd.Blur()
			m.timeStart.Blur()
			m.task.Focus()
			m.step = stepConfirm
		}
	}

	// Update the text inputs
	var cmd tea.Cmd
	m.timeStart, cmd = m.timeStart.Update(msg)
	m.timeEnd, _ = m.timeEnd.Update(msg) // Ignoring second cmd for simplicity

	return m, cmd
}

func (m *Model) validate() {

	startStr := m.timeStart.Value()
	endStr := m.timeEnd.Value()

	if startStr == "" {
		timeStartErr = "Start time cannot be empty."
	}

	if endStr == "" {
		timeEndErr = "End time cannot be empty."
	}

	startTime, err1 := utils.ParseTime(startStr, m.date)
	endTime, err2 := utils.ParseTime(endStr, m.date)

	if err1 != nil {
		timeStartErr = "Invalid start time format. Use HH:MM."
	}

	if err2 != nil {
		timeEndErr = "Invalid end time format. Use HH:MM."
	}

	if err1 == nil && err2 == nil && !endTime.After(startTime) {
		timeEndErr = "End time must be after start time."
	}

}
