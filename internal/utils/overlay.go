package utils

import (
	"clockify-app/internal/styles"
	"strings"

	"charm.land/lipgloss/v2"
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
func RenderWithModal(height, width int, baseContent string, modalContent string) string {
	modalWidth := lipgloss.Width(modalContent)
	modalHeight := lipgloss.Height(modalContent)

	x := (width - modalWidth) / 2
	y := (height - modalHeight) / 2

	base := lipgloss.NewLayer(baseContent).Z(0)
	modal := lipgloss.NewLayer(modalContent).X(x).Y(y).Z(1)

	compositor := lipgloss.NewCompositor(base, modal)

	canvas := lipgloss.NewCanvas(width, height)
	canvas.Compose(compositor)

	return canvas.Render()
}

func RenderScrollbarForModal(contentHeight, maxHeight, offset int) string {
	return RenderScrollbar(contentHeight, maxHeight, offset, 2, styles.CustomBorder.TopRight, styles.CustomBorder.BottomRight)
}

func RenderScrollbarSimple(contentHeight, maxHeight, offset int) string {
	return RenderScrollbar(contentHeight, maxHeight, offset, 0, "↑", "↓")
}

func RenderScrollbar(contentHeight, maxHeight, scrollOffset, borderOffset int, topchar, bottomchar string) string {
	var scrollbar strings.Builder
	viewportHeight := maxHeight + borderOffset
	scrollPercent := float64(scrollOffset) / float64(max(1, contentHeight-maxHeight))

	thumbHeight := max(1, (viewportHeight*maxHeight)/max(1, contentHeight))
	thumbPosition := int(float64(viewportHeight-thumbHeight) * scrollPercent)

	requiresScrollbar := contentHeight > maxHeight

	for i := range viewportHeight {
		var char string
		if i == 0 {
			char = topchar
		} else if i == viewportHeight-1 {
			char = bottomchar
		} else if requiresScrollbar && i-1 >= thumbPosition && i-1 < thumbPosition+thumbHeight {
			char = "█"
		} else {
			char = "│"
		}

		scrollbar.WriteString(lipgloss.NewStyle().Foreground(styles.Secondary).Render(char))
		if i < viewportHeight-1 {
			scrollbar.WriteString("\n")
		}
	}

	return scrollbar.String()
}
