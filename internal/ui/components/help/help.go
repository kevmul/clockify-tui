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

	content := ""

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
