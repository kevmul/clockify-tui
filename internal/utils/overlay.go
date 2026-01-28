package utils

import (
	"clockify-app/internal/styles"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// =======================================
// Get the width of the modal
// based on line width
// =======================================

func WidthOfModal(width int, modal string) int {
	modalLines := strings.Split(modal, "\n")
	modalWidth := 0
	for _, line := range modalLines {
		lineLen := lipgloss.Width(line)
		if lineLen > modalWidth {
			modalWidth = lineLen
		}
	}
	return modalWidth
}

// =======================================
// Custom overlay function
// =======================================

func RenderWithModal(height, width int, baseContent string, modal string) string {
	// Split base content into lines
	baseLines := strings.Split(baseContent, "\n")

	// Ensure we have enough lines for the terminal height
	for len(baseLines) < height {
		baseLines = append(baseLines, "")
	}

	// Render the modal
	modalContent := modal
	modalLines := strings.Split(modalContent, "\n")

	// Calculate starting position to center the modal
	modalHeight := len(modalLines)
	startRow := (height - modalHeight) / 2
	if startRow < 0 {
		startRow = 0
	}

	// Find the actual width of the modal
	modalWidth := WidthOfModal(width, modal)
	startCol := (width - modalWidth) / 2
	if startCol < 0 {
		startCol = 0
	}

	// Helper to truncate string at visual width (ANSI-aware and emoji-aware)
	truncateAt := func(s string, width int) string {
		if width <= 0 {
			return ""
		}
		var result strings.Builder
		currentWidth := 0
		inEscape := false

		for _, r := range s {
			if r == '\x1b' {
				inEscape = true
			}

			if inEscape {
				result.WriteRune(r)
				if r == 'm' {
					inEscape = false
				}
				continue
			}

			if currentWidth >= width {
				break
			}

			runeW := runewidth.RuneWidth(r)
			if currentWidth+runeW > width {
				break
			}

			result.WriteRune(r)
			currentWidth += runeW
		}
		return result.String()
	}

	// Helper to skip first N visual characters (ANSI-aware and emoji-aware)
	skipChars := func(s string, n int) string {
		if n <= 0 {
			return s
		}

		skipped := 0
		inEscape := false
		var result strings.Builder
		started := false

		for _, r := range s {
			if r == '\x1b' {
				inEscape = true
			}

			if started || inEscape {
				result.WriteRune(r)
			}

			if inEscape {
				if r == 'm' {
					inEscape = false
				}
				continue
			}

			if !started {
				runeW := runewidth.RuneWidth(r)
				skipped += runeW
				if skipped > n {
					started = true
					result.WriteRune(r)
				}
			}
		}
		return result.String()
	}

	// Overlay modal lines onto base lines
	for i, modalLine := range modalLines {
		row := startRow + i
		if row >= 0 && row < len(baseLines) {
			baseLine := baseLines[row]
			baseWidth := lipgloss.Width(baseLine)

			// Extract left part (before modal)
			leftPart := truncateAt(baseLine, startCol)

			// Extract right part (after modal)
			endCol := startCol + lipgloss.Width(modalLine)
			var rightPart string
			if endCol < baseWidth {
				rightPart = skipChars(baseLine, endCol)
			}

			// Pad if needed
			leftWidth := lipgloss.Width(leftPart)
			if leftWidth < startCol {
				leftPart += strings.Repeat(" ", startCol-leftWidth)
			}

			baseLines[row] = leftPart + modalLine + rightPart
		}
	}

	return strings.Join(baseLines, "\n")
}

func RenderScrollbarForModal(viewport viewport.Model) string {
	return RenderScrollbar(viewport, "┐", "┘")
}

func RenderScrollbarSimple(viewport viewport.Model) string {
	return RenderScrollbar(viewport, "↑", "↓")
}

func RenderScrollbar(viewport viewport.Model, topchar, bottomchar string) string {
	if viewport.TotalLineCount() <= viewport.Height {
		// No scrollbar needed
		return ""
	}
	totalLines := viewport.TotalLineCount()
	// Add 2 to account for top and bottom borders
	viewportHeight := viewport.Height + 2
	scrollPercent := float64(viewport.YOffset) / float64(totalLines-viewport.Height)

	thumbHeight := max(1, (viewportHeight*viewportHeight)/totalLines)
	thumbPosition := int(float64(viewport.Height-thumbHeight) * scrollPercent)

	var scrollbar strings.Builder
	for i := range viewportHeight {
		var char string
		if i == 0 {
			char = topchar // Top corner
		} else if i == viewportHeight-1 {
			char = bottomchar // Bottom corner
		} else if i-1 >= thumbPosition && i-1 < thumbPosition+thumbHeight {
			char = "█" // Thumb (offset by 1 to account for top corner)
		} else {
			char = "│" // Normal border
		}

		scrollbar.WriteString(lipgloss.NewStyle().Foreground(styles.Primary).Render(char))
		if i < viewportHeight-1 {
			scrollbar.WriteString("\n")
		}
	}

	return scrollbar.String()
}
