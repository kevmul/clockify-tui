package entryform

import (
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
		case "tab":
			// Handle tab to go to next step
			if m.projectSearch.Focused() {
				// If project search is focused, don't move to next step
				break
			}

			if m.step == stepTimeInput {
				// If we're in time input step, ensure end time is focused next
				break
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

			if m.step > stepDateSelect {
				m.step--
			}
		}

	case messages.TasksLoadedMsg:
		m.tasks = msg.Tasks
		m.tasksReady = true
		return m, nil
	}

	switch m.step {
	case stepDateSelect:
		m, cmd = m.updateDateSelect(msg)
	case stepDescriptionInput:
		m, cmd = m.updateDescriptionInput(msg)
	case stepProjectSelect:
		m, cmd = m.updateProjectSelect(msg)
	case stepTimeInput:
		m, cmd = m.updateTimeInput(msg)
	case stepTaskInput:
		m, cmd = m.updateTaskInput(msg)
	case stepConfirm:
		m, cmd = m.updateConfirm(msg)
	case stepComplete:
		m, cmd = m.updateComplete(msg)
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
