package entryform

import (
	"bytes"
	"clockify-app/internal/api"
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// These constants represent each screen in our UI flow
// Using constants (instead of magic numbers) makes the code more readable
const (
	stepDateSelect       = iota //  Select which date to log time for
	stepDescriptionInput        //  Enter task description
	stepProjectSelect           //  Select which project
	stepTaskInput               //  Select a task if applicable
	stepTimeInput               //  Enter time range (e.g., "9a - 5p")
	stepConfirm                 //  Review and confirm the entry
	stepComplete                //  Show success message
)

type Model struct {
	// Current step in the workflow (which screen we're on)
	apiKey      string
	workspaceID string
	step        int
	StepLines   int // Number of lines in the current step's view (for viewport sizing)

	// Data from API
	projects   []models.Project // List of available projects
	tasks      []models.Task    //
	tasksReady bool             // Whether tasks have been loaded

	// Navigation state
	cursor   int // Current position in lists (for arrow key navigation)
	selected int // Index of selected item (not currently used but kept for future)

	// User inputs
	date           time.Time       // Selected date for time entry
	timeStart      textinput.Model // Text input for start time (e.g., "9:00 AM")
	timeEnd        textinput.Model // Text input for end time (e.g., "5:00 PM")
	description    textinput.Model // Text input for task description
	task           textinput.Model // Text input for task description
	projectSearch  textinput.Model // Text input for project search
	selectedProj   models.Project  // The project user selected
	selectedProjID int             // ID of the selected project
	selectedTask   models.Task     // The task user selected
	selectedEntry  models.Entry    // The time entry being edited (if any)

	// Status flags
	editing    bool  // Whether we're in editing mode
	err        error // Any error that occurred
	submitting bool  // Whether we're currently submitting (not used yet)
	success    bool  // Whether submission was successful
}

func New(cfg *config.Config, projects []models.Project) Model {
	// Create and configure the start time input
	timeStartInput := textinput.New()
	timeStartInput.Placeholder = "e.g., 9a"
	timeStartInput.CharLimit = 8 // "12:00 PM" is 8 characters
	timeStartInput.Width = 30

	// Create and configure the time end input
	timeEndInput := textinput.New()
	timeEndInput.Placeholder = "e.g., 9a"
	timeEndInput.CharLimit = 8 // "12:00 PM" is 8 characters
	timeEndInput.Width = 30

	// Create and configure the task name input
	descriptionInput := textinput.New()
	descriptionInput.Placeholder = "Enter task description"
	descriptionInput.CharLimit = 100
	descriptionInput.Width = 50

	// Create and configure the task name input
	taskInput := textinput.New()
	taskInput.Placeholder = "Enter task description"
	taskInput.CharLimit = 100
	taskInput.Width = 50

	// Create and configure the project search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search projects..."
	searchInput.Width = 50

	return Model{
		apiKey:        cfg.APIKey,
		workspaceID:   cfg.WorkspaceId,
		step:          stepDateSelect, // Start at date selection
		date:          time.Now(),     // Default to today
		timeStart:     timeStartInput,
		timeEnd:       timeEndInput,
		description:   descriptionInput,
		task:          taskInput,
		projectSearch: searchInput,
		projects:      projects,
		cursor:        0, // Start at first item in lists
		editing:       false,
		tasksReady:    false,
	}
}

func (m Model) UpdateEntry(entry models.Entry) Model {
	m.editing = true

	m.selectedEntry = entry
	m.date = entry.TimeInterval.Start.In(time.Local)

	m.description.SetValue(entry.Description)

	// Pre-fill time inputs
	startStr := entry.TimeInterval.Start.In(time.Local).Format("3:04 PM")
	endStr := entry.TimeInterval.End.In(time.Local).Format("3:04 PM")
	m.timeStart.SetValue(startStr)
	m.timeEnd.SetValue(endStr)

	// Find and select the project
	for i, proj := range m.projects {
		if proj.ID == entry.ProjectID {
			m.selectedProj = proj
			m.selectedProjID = i
			break
		}
	}

	// Find the cursor position for the project
	for i, proj := range m.projects {
		if proj.ID == entry.ProjectID {
			m.cursor = i
			break
		}
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Set the projects list in the model
func (m Model) SetProjects(projects []models.Project) Model {
	m.projects = projects
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// Implementation of Update method goes here
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global key handling can go here if needed
		switch msg.String() {
		case "esc":
			// Handle escape to exit the form
			// Reset form state if needed
			m = New(&config.Config{APIKey: m.apiKey, WorkspaceId: m.workspaceID}, m.projects)
			timeStartErr = "" // Located in time input step file
			timeEndErr = ""
		case "tab":
			// Handle tab to go to next step
			if m.projectSearch.Focused() {
				// If project search is focused, don't move to next step
				break
			}

			switch m.step {
			case stepTaskInput:
				m.timeStart.Focus()
			case stepTimeInput:
				// If we're in time input step, ensure end time is focused next
				if m.timeStart.Focused() {
					m.timeStart.Blur()
					m.timeEnd.Focus()
				} else {
					m.timeEnd.Blur()
					m.timeStart.Focus()
				}
				_, cmd = m.updateTimeInput(msg)
				return m, nil
				// }
			case stepProjectSelect:
				// If no project selected, don't move forward
				if m.selectedProj.ID == "" {
					return m, nil
				}

			}

			if m.step < stepConfirm {
				m.step++
			}

		case "shift+tab":
			// Handle shift+tab to go back a step
			if m.projectSearch.Focused() {
				// If project search is focused, don't move to next step
				break
			}

			switch m.step {
			case stepProjectSelect:
				// Reset time input errors when going back
				timeStartErr = "" // Located in time input step file
				timeEndErr = ""
				m.description.Focus()
			case stepTimeInput:
				// Blur both inputs
				m.timeStart.Blur()
				m.timeEnd.Blur()
			}

			if m.step > stepDateSelect {
				m.step--
			}
		case "enter":
			switch m.step {

			case stepDateSelect:
				m.step = stepDescriptionInput
				m.description.Focus()

			case stepDescriptionInput:
				m.description.Blur()
				m.step = stepProjectSelect

			case stepProjectSelect:
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
					m.task.Focus()
					m.cursor = 0
					m.tasks = nil
					m.tasksReady = false
					m.step = stepTaskInput
					m.cursor = 0
				}
				return m, api.FetchTasks(m.apiKey, m.workspaceID, m.selectedProj.ID)

			case stepTaskInput:
				if len(m.tasks) > 0 {
					m.selectedTask = m.tasks[m.cursor]
				}
				m.task.Blur()
				m.step = stepTimeInput
				m.timeStart.Focus()

			case stepTimeInput:

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
				m.step = stepConfirm

			case stepConfirm:
				// Submit the time entry and transition to submission state
				m.submitting = true
				m.step++
				if m.editing {
					// Updating an existing entry
					cmds = append(cmds, m.updateTimeEntry())
				} else {
					cmds = append(cmds, m.submitTimeEntry())
				}

			case stepComplete:
				cmds = append(cmds, func() tea.Msg {
					// Reset form after completion
					return messages.ModalClosedMsg{}
				})
			}

		}

	case messages.TasksLoadedMsg:
		m.tasks = msg.Tasks
		m.tasks = append(m.tasks, models.Task{ID: "", Name: "No Task"}) // Option for no task
		m.tasksReady = true
		m.StepLines = getLines(m.viewTimeInput())
		return m, nil
	}

	switch m.step {
	case stepDateSelect:
		m, cmd = m.updateDateSelect(msg)
		m.StepLines = getLines(m.viewDateSelect())
	case stepDescriptionInput:
		m, cmd = m.updateDescriptionInput(msg)
		m.StepLines = getLines(m.viewDescriptionInput())
	case stepProjectSelect:
		m, cmd = m.updateProjectSelect(msg)
		m.StepLines = getLines(m.viewProjectSelect())
	case stepTimeInput:
		m, cmd = m.updateTimeInput(msg)
		m.StepLines = getLines(m.viewTimeInput())
	case stepTaskInput:
		m, cmd = m.updateTaskInput(msg)
	case stepConfirm:
		// m, cmd = m.updateConfirm(msg)
		m.StepLines = getLines(m.viewConfirm())
	case stepComplete:
		m, cmd = m.updateComplete(msg)
		m.StepLines = getLines(m.viewCompletionInput())
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	// Implementation of View method goes here
	s := ""
	switch m.step {
	case stepDateSelect:
		s += m.viewDateSelect()
	case stepDescriptionInput:
		s += m.viewDescriptionInput()
	case stepProjectSelect:
		s += m.viewProjectSelect()
	case stepTimeInput:
		s += m.viewTimeInput()
	case stepTaskInput:
		s += m.viewTaskInput()
	case stepConfirm:
		s += m.viewConfirm()
	case stepComplete:
		s += m.viewCompletionInput()
	default:
		s += "Unknown step"
	}

	return s
}

func getLines(s string) int {
	return bytes.Count([]byte(s), []byte{'\n'})
}
