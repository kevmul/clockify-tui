package help

import "github.com/charmbracelet/bubbles/key"

// =======================================
// Global Key Bindings
// =======================================

type GlobalKeyMap struct {
	Navigation key.Binding
	Help       key.Binding
	Quit       key.Binding
	Up         key.Binding
	Down       key.Binding
	Esc        key.Binding
}

var Global = GlobalKeyMap{
	Navigation: key.NewBinding(
		key.WithKeys("1", "2"),
		key.WithHelp("1 - 0", "Switch view"),
	),
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

// =======================================
// Entry Key Bindings
// =======================================

type EntryKeyMap struct {
	Delete key.Binding
	Edit   key.Binding
	New    key.Binding
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Search key.Binding
}

var Entry = EntryKeyMap{
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "Search entries"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "Move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "Move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "Paginate Left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "Paginate Right"),
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

// =======================================
// Projects Key Bindings
// =======================================

type ProjectsKeyMap struct {
	Enter key.Binding
	Up    key.Binding
	Down  key.Binding
}

var Projects = ProjectsKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<enter>", "Open Project"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "Move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "Move down"),
	),
}

// =======================================
// Project Single Key Bindings
// =======================================

type ProjectKeyMap struct {
	Back key.Binding
}

var Project = ProjectKeyMap{
	Back: key.NewBinding(
		key.WithKeys("b", "esc"),
		key.WithHelp("b/<esc>", "Back to Projects"),
	),
}

// =======================================
// Settings Key Bindings
// =======================================

type SettingsKeyMap struct {
	Enter    key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
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
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<enter>", "Toggle lock / Save"),
	),
}
