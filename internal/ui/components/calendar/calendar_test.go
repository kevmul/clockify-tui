package calendar

import (
	"testing"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func TestNew(t *testing.T) {
	m := New()
	if m.KeyMap.Today.Keys()[0] != "t" {
		t.Error("KeyMap not initialized")
	}
	if m.Styles.Header.GetBold() != true {
		t.Error("Styles not initialized")
	}
}

func TestSetInitalDay(t *testing.T) {
	m := New()
	testDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	m.SetInitalDay(testDate)
	if !m.initialDay.Equal(testDate) {
		t.Errorf("Expected initialDay %v, got %v", testDate, m.initialDay)
	}
}

func TestSetSelectedDay(t *testing.T) {
	m := New()
	testDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	m.SetSelectedDay(testDate)
	if !m.SelectedDate.Equal(testDate) {
		t.Errorf("Expected SelectedDate %v, got %v", testDate, m.SelectedDate)
	}
}

func TestUpdate_NextDay(t *testing.T) {
	m := New()
	m.SelectedDate = time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	m.CurrentDate = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	msg := tea.KeyPressMsg{Code: 'l'}
	m, _ = m.Update(msg)

	expected := time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC)
	if !m.SelectedDate.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, m.SelectedDate)
	}
}

func TestUpdate_PreviousDay(t *testing.T) {
	m := New()
	m.SelectedDate = time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	m.CurrentDate = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	msg := tea.KeyPressMsg{Code: 'h'}
	m, _ = m.Update(msg)

	expected := time.Date(2026, 1, 14, 0, 0, 0, 0, time.UTC)
	if !m.SelectedDate.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, m.SelectedDate)
	}
}

func TestUpdate_NextWeek(t *testing.T) {
	m := New()
	m.SelectedDate = time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	m.CurrentDate = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	msg := tea.KeyPressMsg{Code: 'j'}
	m, _ = m.Update(msg)

	expected := time.Date(2026, 1, 22, 0, 0, 0, 0, time.UTC)
	if !m.SelectedDate.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, m.SelectedDate)
	}
}

func TestUpdate_PreviousWeek(t *testing.T) {
	m := New()
	m.SelectedDate = time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	m.CurrentDate = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	msg := tea.KeyPressMsg{Code: 'k'}
	m, _ = m.Update(msg)

	expected := time.Date(2026, 1, 8, 0, 0, 0, 0, time.UTC)
	if !m.SelectedDate.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, m.SelectedDate)
	}
}

func TestUpdate_NextMonth(t *testing.T) {
	m := New()
	m.SelectedDate = time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	m.CurrentDate = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	msg := tea.KeyPressMsg{Code: tea.KeyPgDown}
	m, _ = m.Update(msg)

	expected := time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC)
	if !m.SelectedDate.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, m.SelectedDate)
	}
	if m.CurrentDate.Month() != time.February {
		t.Errorf("Expected CurrentDate month to be February, got %v", m.CurrentDate.Month())
	}
}

func TestUpdate_PreviousMonth(t *testing.T) {
	m := New()
	m.SelectedDate = time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC)
	m.CurrentDate = time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	msg := tea.KeyPressMsg{Code: tea.KeyPgUp}
	m, _ = m.Update(msg)

	expected := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	if !m.SelectedDate.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, m.SelectedDate)
	}
	if m.CurrentDate.Month() != time.January {
		t.Errorf("Expected CurrentDate month to be January, got %v", m.CurrentDate.Month())
	}
}

func TestUpdate_Today(t *testing.T) {
	m := New()
	m.SelectedDate = time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	msg := tea.KeyPressMsg{Code: 't'}
	m, _ = m.Update(msg)

	now := time.Now()
	if m.SelectedDate.Day() != now.Day() || m.SelectedDate.Month() != now.Month() || m.SelectedDate.Year() != now.Year() {
		t.Errorf("Expected today's date, got %v", m.SelectedDate)
	}
}

func TestUpdate_MonthChange(t *testing.T) {
	m := New()
	m.SelectedDate = time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
	m.CurrentDate = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	msg := tea.KeyPressMsg{Code: 'l'}
	m, _ = m.Update(msg)

	if m.CurrentDate.Month() != time.February {
		t.Errorf("Expected CurrentDate to change to February, got %v", m.CurrentDate.Month())
	}
}

func TestDefaultKeyMap(t *testing.T) {
	km := DefaultKeyMap()
	if !key.Matches(tea.KeyPressMsg{Code: 't'}, km.Today) {
		t.Error("Today key not set correctly")
	}
	if !key.Matches(tea.KeyPressMsg{Code: 'l'}, km.Next) {
		t.Error("Next key not set correctly")
	}
}

func TestView(t *testing.T) {
	m := New()
	m.CurrentDate = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	m.SelectedDate = time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	view := m.View()
	if view.Content == "" {
		t.Error("View should not be empty")
	}
}
