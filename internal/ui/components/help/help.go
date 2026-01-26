package help

import (
	"clockify-app/internal/styles"
	"reflect"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	sections []HelpSection
	width    int
	height   int
}

// HelpSection represents a section of help bindings
type HelpSection struct {
	Title   string
	Binding []key.Binding
}

func New(sections ...HelpSection) Model {
	return Model{
		sections: sections,
		width:    80,
		height:   24,
	}
}

func (m *Model) setSize(width, height int) {
	m.width = width
	m.height = height
}

func (m Model) View() string {

	content := "Key Bindings \n\n"

	bindingStyle := lipgloss.NewStyle()

	for i, section := range m.sections {

		content += styles.TitleStyle.MarginBottom(0).Render(section.Title) + "\n"

		// Bindings
		for _, binding := range section.Binding {
			if !binding.Enabled() {
				continue
			}

			keyStr := binding.Help().Key
			desc := binding.Help().Desc

			line := lipgloss.JoinHorizontal(
				lipgloss.Top,
				styles.MutedTextStyle.Width(12).Align(lipgloss.Right).MarginRight(2).Render(keyStr),
				desc,
			)
			content += bindingStyle.Render(line) + "\n"
		}
		if i < len(m.sections)-1 {
			content += "\n"
		}
	}

	return content
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg any) (Model, tea.Cmd) {
	return m, nil
}

func GenerateSection(title string, km any) HelpSection {
	return HelpSection{
		Title:   title,
		Binding: bindingsFromStruct(km),
	}
}

type GlobalKeyMap struct {
	Help key.Binding
	Quit key.Binding
	Up   key.Binding
	Down key.Binding
	Esc  key.Binding
}

var Global = GlobalKeyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Show help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/<ctrl+c>", "Quit"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "Move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "Move down"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("<esc>", "Close help"),
	),
}

type EntryKeyMap struct {
	Delete key.Binding
	Edit   key.Binding
	New    key.Binding
	Search key.Binding
}

var Entry = EntryKeyMap{
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "Search entries"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "New entry"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "Edit entry"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete entry"),
	),
}

type SettingsKeyMap struct {
	LockToggle key.Binding
	Tab        key.Binding
	ShiftTab   key.Binding
}

var Settings = SettingsKeyMap{
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("<tab>", "Next Input"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("<shift+tab>", "Previous Input"),
	),
	LockToggle: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<enter>", "Toggle lock"),
	),
}

// bindingsFromStruct extracts all key.Binding fields from a struct
func bindingsFromStruct(km interface{}) []key.Binding {
	var bindings []key.Binding

	v := reflect.ValueOf(km)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Type() == reflect.TypeOf(key.Binding{}) {
			bindings = append(bindings, field.Interface().(key.Binding))
		}
	}

	return bindings
}
