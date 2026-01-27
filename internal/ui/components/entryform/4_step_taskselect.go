package entryform

import (
	"clockify-app/internal/styles"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ================ Task Selection =================
func (m Model) viewTaskInput() string {
	if !m.tasksReady {
		return styles.SubtitleStyle.Render("Loading tasks...")
	}

	// Implementation of task input view goes here
	title := styles.TitleStyle.MarginBottom(0).Render("Select Task")
	subtitle := styles.SubtitleStyle.MarginBottom(1).Render("Describe the work you did during this time period.")

	sb := strings.Builder{}
	sb.WriteString(title + "\n")
	sb.WriteString(subtitle + "\n")

	if m.tasksReady && len(m.tasks) == 0 {
		sb.WriteString("No tasks found in this workspace.\n\n")
		sb.WriteString(styles.HelpStyle.Render("Press Enter to continue, or Tab/Shift+Tab to navigate."))
		return sb.String()
	}

	const visibleItems = 5 // Show 5 items at a time
	start := 0
	end := len(m.tasks)

	filteredTasks := m.tasks

	// If we have more tasks than can fit, show a window around cursor
	if len(filteredTasks) > visibleItems {
		// Center the cursor in the window
		start = m.cursor - visibleItems/2
		end = start + visibleItems

		// Adjust if we're near the beginning
		if start <= 0 {
			start = 0
			end = visibleItems + 1
		}

		// Adjust if we're near the end
		if end > len(filteredTasks) {
			end = len(filteredTasks)
			start = end - visibleItems + 1
			if start < 0 {
				start = 0
			}
		}

		// Show indicator if there are items above
		if start > 0 {
			sb.WriteString("\n  ↑ More tasks above...\n")
		}
	}

	// List visible tasks
	for i := start; i < end; i++ {
		task := filteredTasks[i]
		if m.cursor == i {
			sb.WriteString(styles.SelectedItemStyle.Render(fmt.Sprintf("❯ %s", task.Name)) + "\n")
		} else {
			sb.WriteString(fmt.Sprintf("  %s\n", task.Name))
		}
	}

	// Show indicator if there are items below
	if end < len(filteredTasks) {
		sb.WriteString("\n  ↓ More tasks below...\n")
	}

	sb.WriteString("\n\n" + styles.HelpStyle.Render("Use arrow keys to navigate, Enter to select."))

	return sb.String()
}

func (m Model) updateTaskInput(msg tea.Msg) (Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}
		case "enter":
			// Move to next step or finish
			if len(m.tasks) > 0 {
				m.selectedTask = m.tasks[m.cursor]
			}
			m.task.Blur()
			m.timeStart.Focus()
			m.step = stepTimeInput
		}
	}
	m.task, cmd = m.task.Update(msg)

	return m, cmd
}
