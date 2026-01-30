package ui

import (
	"clockify-app/internal/api"
	"clockify-app/internal/cache"
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"clockify-app/internal/styles"
	"clockify-app/internal/utils"
	debug "clockify-app/internal/utils"
	"strconv"
	"strings"

	"clockify-app/internal/ui/components/help"
	"clockify-app/internal/ui/components/modal"
	"clockify-app/internal/ui/views/entries"
	"clockify-app/internal/ui/views/project"
	"clockify-app/internal/ui/views/projects"
	"clockify-app/internal/ui/views/settings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type View int

const (
	SettingsView View = iota
	EntriesView
	ProjectsView
	ProjectView
)

type Page struct {
	Label string
	Key   View
}

var pages = []Page{
	{"Entries", EntriesView},
	{"Projects", ProjectsView},
	{"Settings", SettingsView},
}

type Model struct {
	// Config and shared state
	config      *config.Config
	userId      string
	workspaceId string
	projects    []models.Project

	// Current View
	currentView View

	// View models
	settingsView settings.Model // Settings View
	entriesView  entries.Model  // List of Entries
	projectsView projects.Model // List of Projects
	projectView  project.Model  // Single Project view

	// Modal state
	modal     *modal.Model
	showModal bool

	// UI Dimensions
	width  int
	height int

	// Loading state
	ready    bool
	viewport viewport.Model
}

func NewModel() Model {
	cfg, _ := config.LoadConfig()

	// Start at settings if no config
	currentView := SettingsView

	if cfg.APIKey != "" && cfg.WorkspaceId != "" {
		currentView = EntriesView
	}

	return Model{
		config:       cfg,
		currentView:  currentView,
		settingsView: settings.New(cfg),
		entriesView:  entries.New(cfg),
		projectsView: projects.New(cfg),
		ready:        false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.initializeFirstViewCmd(),
		tea.EnterAltScreen,
	)
}

func (m Model) initializeFirstViewCmd() tea.Cmd {
	switch m.currentView {
	case SettingsView:
		return settings.Init()
	case EntriesView:
		return tea.Sequence(
			api.FetchProjects(
				m.config.APIKey,
				m.config.WorkspaceId,
			),
			m.entriesView.Init(),
		)
	case ProjectsView:
		return m.projectsView.Init()
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := 4
		footerHeight := 1
		verticalMarginHeight := headerHeight + footerHeight

		// Update viewport size
		if !m.ready {
			m.viewport = viewport.New(msg.Width-1, msg.Height-verticalMarginHeight)
			// m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.renderContent())
			m.width = msg.Width
			m.height = msg.Height
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 1
			m.viewport.Height = msg.Height - verticalMarginHeight
			m.width = msg.Width
			m.height = msg.Height
		}

		m.projectsView, cmd = m.projectsView.Update(msg)
		cmds = append(cmds, cmd)

		m.entriesView, cmd = m.entriesView.Update(msg)
		cmds = append(cmds, cmd)

		m.projectsView, cmd = m.projectsView.Update(msg)
		cmds = append(cmds, cmd)

		m.settingsView, cmd = m.settingsView.Update(msg)
		cmds = append(cmds, cmd)

		return m, nil

	case tea.KeyMsg:
		// Global keybindings
		if !m.showModal {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "1", "2", "3":
				if num, err := strconv.Atoi(msg.String()); err == nil {
					m.currentView = pages[num-1].Key
					m.viewport.SetContent(m.renderContent())
					// Initialize view if needed
					switch m.currentView {
					case EntriesView:
						return m, tea.Sequence(
							api.FetchProjects(
								m.config.APIKey,
								m.config.WorkspaceId,
							),
							m.entriesView.Init(),
						)
					case ProjectsView:
						return m, m.projectsView.Init()
					case SettingsView:
						return m, settings.Init()
					}
					return m, nil
				}
			case "n":
				switch m.currentView {
				case EntriesView:
					m.showModal = true
					m.modal = modal.NewEntryForm(m.config, m.projects)
					return m, nil
				}
			case "?":
				m.showModal = true
				switch m.currentView {
				case EntriesView:
					m.modal = modal.NewHelp(
						help.GenerateSection("Entries Keys", help.Entry),
						help.GenerateSection("Global Keys", help.Global),
					)
					return m, nil
				case ProjectsView:
					m.modal = modal.NewHelp(
						help.GenerateSection("Projects Keys", help.Projects),
						help.GenerateSection("Global Keys", help.Global),
					)
					return m, nil
				case ProjectView:
					m.modal = modal.NewHelp(
						help.GenerateSection("Project Keys", help.Project),
						help.GenerateSection("Global Keys", help.Global),
					)
					return m, nil
				case SettingsView:
					m.modal = modal.NewHelp(
						help.GenerateSection("Settings Keys", help.Settings),
						help.GenerateSection("Global Keys", help.Global),
					)
					return m, nil
				default:
					m.modal = modal.NewHelp(
						help.GenerateSection("Global Keys", help.Global),
					)
				}
			}
		}

	case messages.UserLoadedMsg:
		m.userId = msg.UserId
		m.settingsView, cmd = m.settingsView.Update(msg)
		m.viewport.SetContent(m.renderContent())
		return m, cmd

	case messages.ConfigSavedMsg:
		m.config = msg.Config
		m.userId = msg.UserId
		m.workspaceId = msg.WorkspaceId
		m.viewport.SetContent(m.renderContent())
		_ = m.config.Save()
		return m, nil

	case messages.ProjectsLoadedMsg:
		m.projects = msg.Projects
		switch m.currentView {
		case EntriesView:
			// Let entries view handle the loaded projects
			m.entriesView, cmd = m.entriesView.Update(msg)
		case ProjectsView:
			m.projectsView, cmd = m.projectsView.Update(msg)
		}
		m.viewport.SetContent(m.renderContent())
		return m, cmd

	case messages.ProjectSelectedMsg:
		// Let projects view handle the selected project
		m.projectView = project.New(m.config, msg.Project, []models.Task{})
		cmd = m.projectView.Init()
		m.currentView = ProjectView
		// m.projectView, cmd = m.projectView.Update(msg)
		m.viewport.SetContent(m.renderContent())
		return m, cmd

	case messages.TasksLoadedMsg:
		// Let project view handle the loaded tasks
		debug.Log("Tasks loaded message received in app.go")
		if m.currentView == ProjectView {
			m.projectView, cmd = m.projectView.Update(msg)
			m.viewport.SetContent(m.renderContent())
			return m, cmd
		}

	case messages.ExitViewMsg:
		// Go back to Projects View
		if m.currentView == ProjectView {
			m.currentView = ProjectsView
			m.viewport.SetContent(m.renderContent())
		}
		return m, nil

	case messages.EntrySavedMsg:
		m.showModal = false
		cache := cache.GetInstance()
		cache.AddEntry(msg.Entry)
		m.entriesView, cmd = m.entriesView.Update(msg)
		return m, api.FetchEntries(
			m.config.APIKey,
			m.config.WorkspaceId,
			m.config.UserId,
		)

	case messages.EntryUpdatedMsg:
		m.showModal = false
		cache := cache.GetInstance()
		cache.UpdateEntry(msg.Entry)
		m.entriesView, cmd = m.entriesView.Update(msg)
		return m, api.FetchEntries(
			m.config.APIKey,
			m.config.WorkspaceId,
			m.config.UserId,
		)

	case messages.EntriesLoadedMsg:
		m.entriesView, cmd = m.entriesView.Update(msg)
		m.viewport.SetContent(m.renderContent())
		return m, cmd

	case messages.EntryUpdateStartedMsg:
		m.showModal = true
		m.modal = modal.UpdateEntryForm(m.config, m.projects, msg.Entry)
		m.viewport.SetContent(m.renderContent())
		return m, nil

	case messages.EntryDeleteStartedMsg:
		m.showModal = true
		m.modal = modal.NewDeleteConfirmation(msg.EntryId)
		m.viewport.SetContent(m.renderContent())
		return m, nil

	case messages.ModalClosedMsg:
		m.showModal = false
		m.viewport.SetContent(m.renderContent())
		return m, nil

	case messages.ItemDeletedMsg:
		m.showModal = false
		switch msg.Type {
		case "entry":
			// Let entries view handle the deletion
			cache := cache.GetInstance()
			cache.DeleteEntry(msg.ID)
			m.entriesView, cmd = m.entriesView.Update(msg)
			return m, cmd
		}
	}

	// Route to modal if showing
	if m.showModal && m.modal != nil {
		*m.modal, cmd = m.modal.Update(msg)
		return m, cmd
	}

	// Route to active view
	switch m.currentView {
	case SettingsView:
		m.settingsView, cmd = m.settingsView.Update(msg)
	case ProjectsView:
		m.projectsView, cmd = m.projectsView.Update(msg)
	case EntriesView:
		m.entriesView, cmd = m.entriesView.Update(msg)
	case ProjectView:
		m.projectView, cmd = m.projectView.Update(msg)
	}
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	if _, ok := msg.(tea.KeyMsg); ok {
		// Update viewport content on key events
		m.viewport.SetContent(m.renderContent())
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {

	if !m.ready {
		return "Loading..."
	}

	// Create a top navigation bar

	// Tabs
	navBar := m.RenderNavBar("entries", m.width)

	// Add scrollbar
	scrollbar := utils.RenderScrollbarSimple(m.viewport)

	switch m.currentView {
	case EntriesView, ProjectsView:
		scrollbar = ""
	}
	// The viewport already contains the view content in Update
	viewportView := m.viewport.View()

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		viewportView,
		scrollbar,
	)

	// Overlay modal if showing
	if m.showModal && m.modal != nil {
		content = utils.RenderWithModal(m.height-5, m.width, content, m.modal.View())
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		navBar,
		content,
		styles.InfoBarStyle.Width(m.width).Render("[?]: help, [q][ctrl+c]: quit"),
	)
}

func (m Model) Shutdown() tea.Cmd {
	return tea.ExitAltScreen
}

// Example helper function to render a tab
func RenderTab(label, key string, isActive bool) string {
	keyStyle := lipgloss.NewStyle().Foreground(styles.Muted)

	if isActive {
		keyStyle := lipgloss.NewStyle().Foreground(styles.Primary)
		return styles.ActiveTabStyle.Render(keyStyle.Render(key) + " " + label)
	}

	return styles.InactiveTabStyle.Render(keyStyle.Render(key) + " " + label)
}

// Example helper function to render the full nav bar
func (m Model) RenderNavBar(activeTab string, docWidth int) string {

	tabs := []string{}

	for i, page := range pages {
		if page.Key == m.currentView {
			activeTab = page.Label
		}
		tab := RenderTab(page.Label, string(rune('1'+i)), page.Label == activeTab)
		tabs = append(tabs, tab)
	}

	sep := styles.SeparatorStyle.Render("|")

	fullNav := lipgloss.JoinHorizontal(
		lipgloss.Center,
		strings.Join(tabs, sep),
	)

	return styles.NavContainerStyle.Render(fullNav)
}

func (m Model) renderContent() string {
	// Render active view
	switch m.currentView {
	case SettingsView:
		return m.settingsView.View()
	case ProjectsView:
		return m.projectsView.View()
	case EntriesView:
		return m.entriesView.View()
	case ProjectView:
		return m.projectView.View()
	}

	return ""
}
