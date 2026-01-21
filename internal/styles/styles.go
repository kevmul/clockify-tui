package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	Primary   = lipgloss.Color("#7C3AED")
	Secondary = lipgloss.Color("#EC4899")
	Success   = lipgloss.Color("#10B981")
	Error     = lipgloss.Color("#EF4444")
	Warning   = lipgloss.Color("#F59E0B")
	Muted     = lipgloss.Color("#6B7280")

	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success)

	InfoStyle = lipgloss.NewStyle().
			Foreground(Secondary)

	// Box styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	ModalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2).
			Width(60)

	// Input styles
	FocusedInputStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(Primary).
				Padding(0, 1)

	BlurredInputStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(Muted).
				Padding(0, 1)

	// List styles
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Bold(true).
				PaddingLeft(2)

	NormalItemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	// Help styles
	HelpStyle = lipgloss.NewStyle().
			Foreground(Muted).
			MarginTop(1)

	KeyStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)
)
