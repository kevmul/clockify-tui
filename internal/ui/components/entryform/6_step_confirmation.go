package entryform

import (
	"clockify-app/internal/api"
	"clockify-app/internal/messages"
	"clockify-app/internal/styles"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ================ Confirmation Selection =================
func (m Model) viewConfirm() string {
	// Implementation of confirmation view goes here
	chosenDate := m.date.Format("January 2, 2006")
	chosenStart := m.timeStart.Value()
	chosenEnd := m.timeEnd.Value()
	chosenDescription := m.description.Value()
	chosenTask := m.selectedTask.Name
	choseProject := m.selectedProj

	confirmationBtn := styles.ActiveButtonStyle
	confirmationBtnText := "Create"
	confirmationText := "Press Enter to confirm, or Tab/Shift+Tab to navigate."

	if m.editing {
		confirmationBtnText = "Update"
		confirmationText = "Press Enter to update, or Tab/Shift+Tab to navigate."
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		styles.TitleStyle.Margin(0, 0).Render("Confirm Time Entry"),
		styles.SubtitleStyle.Margin(0, 0, 1, 0).Render("Please review your time entry details:"),
		fmt.Sprintf("📅 Date: %s", chosenDate),
		fmt.Sprintf("⏰ Time: %s - %s", chosenStart, chosenEnd),
		fmt.Sprintf("📝 Description: %s", chosenDescription),
		fmt.Sprintf("📁 Project: %s (%s)", choseProject.Name, choseProject.ClientName),
		fmt.Sprintf("🗂️ Task: %s\n", chosenTask),

		confirmationBtn.Render(confirmationBtnText),
		styles.HelpStyle.Render(confirmationText),
	)
}

// func (m Model) updateConfirm(msg tea.Msg) (Model, tea.Cmd) {
// 	// Implementation of confirmation update goes here
// 	return m, nil
// }

// submitTimeEntry creates a command to submit the time entry
func (m Model) submitTimeEntry() tea.Cmd {
	return createTimeEntry(
		m.apiKey,
		m.workspaceID,
		m.selectedProj.ID,
		m.selectedTask.ID,
		m.description.Value(),
		m.timeStart.Value(),
		m.timeEnd.Value(),
		m.date,
	)
}

func (m Model) updateTimeEntry() tea.Cmd {
	return updateTimeEntry(
		m.apiKey,
		m.workspaceID,
		m.selectedEntry.ID,
		m.selectedProj.ID,
		m.selectedTask.ID,
		m.description.Value(),
		m.timeStart.Value(),
		m.timeEnd.Value(),
		m.date,
	)
}

// createTimeEntry returns a command that creates a time entry
// When complete, it sends either submitSuccessMsg or errMsg
func createTimeEntry(apiKey, workspaceID, projectID, taskID, description, startTime, endTime string, date time.Time) tea.Cmd {
	return func() tea.Msg {
		client := api.NewClient(apiKey)
		entry, err := client.CreateTimeEntry(workspaceID, projectID, taskID, description, startTime, endTime, date)

		if err != nil {
			return messages.ErrorMsg{Err: err}
		}

		// Success - return success message
		return messages.EntrySavedMsg{
			Entry: entry,
		}
	}
}

func updateTimeEntry(apiKey, workspaceID, entryID, projectID, taskID, description, startTime, endTime string, date time.Time) tea.Cmd {
	return func() tea.Msg {
		client := api.NewClient(apiKey)
		entry, err := client.UpdateTimeEntry(workspaceID, entryID, projectID, taskID, description, startTime, endTime, date)

		if err != nil {
			return messages.ErrorMsg{Err: err}
		}

		// Success - return success message
		return messages.EntryUpdatedMsg{
			Entry: entry,
		}
	}
}
