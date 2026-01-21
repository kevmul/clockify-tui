package settings

import (
	"clockify-app/internal/api"
	"clockify-app/internal/config"
	"clockify-app/internal/messages"
	"clockify-app/internal/models"
	"clockify-app/internal/styles"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focusIndex int

const (
	apiKeyInput focusIndex = iota
	workspaceInput
	saveButton
)

type Model struct {
	config         *config.Config
	apiKeyInput    textinput.Model
	workspaceInput textinput.Model
	workspaces     []models.Workspace

	currentIndex            focusIndex
	saving                  bool
	saved                   bool
	err                     error
	userId                  string
	selectedWorkespaceIndex int
	showWorkspacesList      bool
}

func New(cfg *config.Config) Model {
	apiKey := textinput.New()
	apiKey.Placeholder = "Enter your Clockify API Key"
	apiKey.Focus()
	apiKey.CharLimit = 64
	apiKey.Width = 50
	apiKey.EchoMode = textinput.EchoPassword
	apiKey.EchoCharacter = '•'

	if cfg.APIKey != "" {
		apiKey.SetValue(cfg.APIKey)
		apiKey.Blur()
	}

	workspace := textinput.New()
	workspace.Placeholder = "Select your Clockify Workspace"
	workspace.CharLimit = 64
	workspace.Width = 50

	if cfg.WorkspaceId != "" {
		workspace.SetValue(cfg.WorkspaceName)
		workspace.Blur()
	}

	return Model{
		config:         cfg,
		apiKeyInput:    apiKey,
		workspaceInput: workspace,
		currentIndex:   apiKeyInput,
		workspaces:     []models.Workspace{},
	}
}

// Init is the initial command for the settings model.
func Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			m.saved = false
			m.err = nil

			s := msg.String()
			if s == "up" || s == "shift+tab" {
				m.currentIndex--
			} else {
				m.currentIndex++
			}

			if m.currentIndex > saveButton {
				m.currentIndex = apiKeyInput
			} else if m.currentIndex < apiKeyInput {
				m.currentIndex = saveButton
			}

			return m, m.updateFocus()

		case "enter":
			switch m.currentIndex {
			case workspaceInput:
				// Show workspace list if we have an API key
				if m.apiKeyInput.Value() != "" {
					m.showWorkspacesList = true
					return m, m.fetchWorkspaces()
				}

			case saveButton:
				return m, m.saveConfig()
			}

		case "esc":
			m.showWorkspacesList = false
			m.selectedWorkespaceIndex = 0

		default:
			// Navigate workspace list
			if m.showWorkspacesList {
				switch msg.String() {
				case "j", "down":
					if m.selectedWorkespaceIndex < len(m.workspaces)-1 {
						m.selectedWorkespaceIndex++
					}
				case "k", "up":
					if m.selectedWorkespaceIndex > 0 {
						m.selectedWorkespaceIndex--
					}
				case "enter":
					if len(m.workspaces) > 0 {
						selectedWorkspace := m.workspaces[m.selectedWorkespaceIndex]
						m.workspaceInput.SetValue(selectedWorkspace.Name)
						m.showWorkspacesList = false
						m.selectedWorkespaceIndex = 0
					}
					return m, nil
				}
			}
		}

	case messages.WorkspacesLoadedMsg:
		m.workspaces = msg.Workspaces
		m.saving = false
		return m, nil

	case messages.ConfigSavedMsg:
		m.saving = false
		m.saved = true
		m.config = msg.Config
		m.userId = msg.UserId
		return m, nil

	case messages.ErrorMsg:
		m.saving = false
		m.err = msg.Err
		return m, nil
	}

	// Update active input
	if !m.showWorkspacesList {
		switch m.currentIndex {
		case apiKeyInput:
			m.apiKeyInput, cmd = m.apiKeyInput.Update(msg)
		case workspaceInput:
			m.workspaceInput, cmd = m.workspaceInput.Update(msg)
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("⚙️ Settings"))
	b.WriteString("\n\n")

	// API Key Input
	b.WriteString(m.renderLabel("Clockify API Key:", apiKeyInput))
	b.WriteString("\n")
	b.WriteString(m.renderInput(m.apiKeyInput, apiKeyInput))
	b.WriteString("\n")
	b.WriteString(styles.SubtitleStyle.Render(" Get your API key from https://clockify.me/user/settings"))
	b.WriteString("\n\n")

	// Workspace Input
	b.WriteString(m.renderLabel("Clockify Workspace:", workspaceInput))
	b.WriteString("\n")
	b.WriteString(m.renderInput(m.workspaceInput, workspaceInput))
	b.WriteString("\n")
	b.WriteString(styles.SubtitleStyle.Render(" Select your workspace (requires valid API key)"))
	b.WriteString("\n\n")

	// Show workspaces list if active
	if m.showWorkspacesList {
		b.WriteString(m.renderWorkspacesList())
		b.WriteString("\n\n")
	}

	// Current config info
	if m.config.UserId != "" {
		b.WriteString(styles.SuccessStyle.Render(fmt.Sprintf("✓ Current User ID: %s", m.config.UserId)))
		b.WriteString("\n")
	}
	if m.config.WorkspaceId != "" {
		b.WriteString(styles.SuccessStyle.Render(fmt.Sprintf("✓ Current Workspace ID: %s", m.config.WorkspaceId)))
		b.WriteString("\n\n")
	}

	// Save Button
	b.WriteString(m.renderSaveButton())
	b.WriteString("\n\n")

	// Status Messages
	if m.saving {
		b.WriteString(styles.InfoStyle.Render("Saving configuration..."))
		b.WriteString("\n")
	} else if m.saved {
		b.WriteString(styles.SuccessStyle.Render("✓ Configuration saved successfully!"))
		b.WriteString("\n")
	} else if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("✗ Error: %s", m.err.Error())))
		b.WriteString("\n")
	}

	// Help
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("Tab/Shift+Tab: Navigate • Enter: Select/Save • Esc: Close List • Ctrl+C: Quit"))

	return lipgloss.NewStyle().Padding(1, 2).Render(b.String())
}

// Helper to render input labels with focus style
func (m Model) renderLabel(label string, index focusIndex) string {
	style := lipgloss.NewStyle()
	if m.currentIndex == index {
		style = style.Foreground(styles.Primary).Bold(true)
	} else {
		style = style.Foreground(styles.Muted)
	}
	return style.Render(label)
}

// Helper to render inputs with focus style
func (m Model) renderInput(input textinput.Model, index focusIndex) string {
	if m.currentIndex == index {
		return styles.FocusedInputStyle.Render(input.View())
	}
	return styles.BlurredInputStyle.Render(input.View())
}

// Helper to render the save button with focus style
func (m Model) renderSaveButton() string {
	label := " Save Configuration "
	// style := lipgloss.NewStyle().Padding(0, 2)
	if m.currentIndex == saveButton {
		// style = style.Background(ui.Primary).Foreground(ui.Background).Bold(true)
		return styles.FocusedInputStyle.
			Width(len(label)).
			Align(lipgloss.Center).
			Render(label)
	}
	// style = style.Background(ui.Muted).Foreground(ui.Background)
	return styles.BlurredInputStyle.
		Width(len(label)).
		Align(lipgloss.Center).
		Render(label)
}

// Helper to render the workspaces list
func (m Model) renderWorkspacesList() string {
	if len(m.workspaces) == 0 {
		return styles.SubtitleStyle.Render(" Loading workspaces...")
	}

	var b strings.Builder
	b.WriteString(styles.BoxStyle.Width(50).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.renderWorkspaceItems()...,
		),
	))

	return b.String()
}

// Helper to render individual workspace items
func (m Model) renderWorkspaceItems() []string {
	items := []string{
		styles.TitleStyle.Render(" Select a Workspace "),
		"",
	}

	for i, ws := range m.workspaces {
		cursor := "  "
		style := styles.NormalItemStyle

		if i == m.selectedWorkespaceIndex {
			cursor = "➤ "
			style = styles.SelectedItemStyle
		}

		items = append(items, style.Render(cursor+ws.Name))
	}

	items = append(items, "", styles.SubtitleStyle.Render("j/k or ↑/↓: navigate • enter: select • esc: cancel"))

	return items
}

// Helper to update focus styles
func (m Model) updateFocus() tea.Cmd {
	cmds := []tea.Cmd{}

	switch m.currentIndex {
	case apiKeyInput:
		cmds = append(cmds, m.apiKeyInput.Focus())
		m.workspaceInput.Blur()
	case workspaceInput:
		cmds = append(cmds, m.workspaceInput.Focus())
		m.apiKeyInput.Blur()
	default:
		m.apiKeyInput.Blur()
		m.workspaceInput.Blur()
	}

	return tea.Batch(cmds...)
}

// Helper to fetch workspaces based on API key
func (m Model) fetchWorkspaces() tea.Cmd {
	return func() tea.Msg {
		client := api.NewClient(m.apiKeyInput.Value())
		workspaces, err := client.GetWorkspaces()
		if err != nil {
			return messages.ErrorMsg{Err: err}
		}
		return messages.WorkspacesLoadedMsg{Workspaces: workspaces}
	}
}

// Helper to save the configuration
func (m Model) saveConfig() tea.Cmd {
	return func() tea.Msg {
		m.saving = true
		return messages.ConfigSavedMsg{
			Config: m.config,
			UserId: m.userId,
		}
	}
}
