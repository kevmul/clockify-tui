package utils

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

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
	modalWidth := 0
	for _, line := range modalLines {
		lineLen := lipgloss.Width(line)
		if lineLen > modalWidth {
			modalWidth = lineLen
		}
	}

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
