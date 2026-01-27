package confirmation

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	model := New("entry1", "entry")

	if model.itemToDelete != "entry1" {
		t.Errorf("Expected itemToDelete 'entry1', got %q", model.itemToDelete)
	}
	if model.itemType != "entry" {
		t.Errorf("Expected itemType 'entry', got %q", model.itemType)
	}
	if model.cursor != 0 {
		t.Errorf("Expected cursor to be 0 initially, got %d", model.cursor)
	}
}

func TestUpdate(t *testing.T) {
	model := New("entry1", "entry")

	// Test left/right arrow keys for cursor movement
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	if updated.cursor != 1 {
		t.Errorf("Expected cursor to be 1 after right arrow, got %d", updated.cursor)
	}

	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if updated.cursor != 0 {
		t.Errorf("Expected cursor to be 0 after left arrow, got %d", updated.cursor)
	}
}
