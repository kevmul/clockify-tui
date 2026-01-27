package entryform

import (
	"clockify-app/internal/models"
	"clockify-app/internal/styles"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ================ Project Selection =================
func (m Model) viewProjectSelect() string {
	// Implementation of project selection view goes here
	// title := styles.TitleStyle.Margin(0, 0).Render("Select Project")
	// subtitle := styles.SubtitleStyle.Margin(0, 0, 1, 0).Render("Use arrow keys to navigate, Enter to select")

	sb := strings.Builder{}

	// Title and subtitle
	sb.WriteString(styles.TitleStyle.Margin(0, 0).Render("Select Project") + "\n")
	sb.WriteString(styles.SubtitleStyle.Margin(0, 0, 1, 0).Render("Use arrow keys to navigate, Enter to select") + "\n")

	// Show search input
	sb.WriteString("🔍 " + m.projectSearch.View() + "\n\n")

	// Filter projects based on search
	filteredProjects := m.filterProjects()

	if len(filteredProjects) == 0 {
		sb.WriteString("  No projects match your search.\n\n")
	}

	// Calculate visible range for scrolling
	const visibleItems = 5 // Show 10 items at a time
	start := 0
	end := len(filteredProjects)

	// If we have more projects than can fit, show a window around cursor
	if len(filteredProjects) > visibleItems {
		// Center the cursor in the window
		start = m.cursor - visibleItems/2
		end = start + visibleItems

		// Adjust if we're near the beginning
		if start <= 0 {
			start = 0
			end = visibleItems + 1
		}

		// Adjust if we're near the end
		if end > len(filteredProjects) {
			end = len(filteredProjects)
			start = end - visibleItems + 1
			if start < 0 {
				start = 0
			}
		}

		// Show indicator if there are items above
		if start > 0 {
			sb.WriteString(fmt.Sprintf("  ↑ %d more above...\n", start))
		}
	}

	// Show visible projects
	for i := start; i < end; i++ {
		proj := filteredProjects[i]

		// Format project name with client if available
		displayName := fmt.Sprintf("%s", proj.Name)
		if proj.ClientName != "" {
			displayName = fmt.Sprintf("%s (%s)", proj.Name, proj.ClientName)
		}

		if m.cursor == i {
			// This is the selected item -= highlight it
			sb.WriteString(styles.SelectedItemStyle.Render(fmt.Sprintf("❯ %s", displayName)) + "\n")
		} else {
			// Regular rendering for unselected items
			sb.WriteString(fmt.Sprintf("  %s\n", displayName))
		}
	}

	// Show indicator if there are items below
	if len(filteredProjects) > visibleItems && end < len(filteredProjects) {
		sb.WriteString(fmt.Sprintf("  ↓ %d more below...", len(filteredProjects)-end))
	}

	return sb.String()
}

// updateProjectSelect handles messages for the project selection step.
func (m Model) updateProjectSelect(msg tea.Msg) (Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.filterProjects())-1 {
				m.cursor++
			}

		case "/":
			// Focus the search input
			if m.projectSearch.Focused() {
				m.projectSearch.Blur()
				return m, nil
			}

			return m, m.projectSearch.Focus()

		case "enter":
			// Select the current project
			if m.projectSearch.Focused() {
				// If search is focused, do nothing on enter
				if m.projectSearch.Focused() {
					m.projectSearch.Blur()
					m.cursor = 0 // Reset cursor when focusing search
					return m, nil
				}
				return m, nil
			}
			filtered := m.filterProjects()
			if len(filtered) > 0 && m.cursor < len(filtered) {
				m.selectedProj = filtered[m.cursor]
				// m.selectedProj = filteredProjects[m.cursor] // Save selected project
				// Move to next step
				m.timeStart.Focus()
				m.step++
				// Reset cursor for next step
				m.cursor = 0
			}
		}
	}

	// Update the project search input
	m.projectSearch, cmd = m.projectSearch.Update(msg)
	cmds = append(cmds, cmd)

	return m, nil
}

// filterProjects filters the list of projects based on the current search query.
func (m Model) filterProjects() []models.Project {
	query := strings.ToLower(strings.TrimSpace(m.projectSearch.Value()))
	if query == "" {
		return m.projects
	}

	var filtered []models.Project
	for _, proj := range m.projects {
		if strings.Contains(strings.ToLower(proj.Name), query) {
			filtered = append(filtered, proj)
		}
	}
	return filtered
}
