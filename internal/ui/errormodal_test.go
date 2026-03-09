package ui

import (
	"testing"

	"github.com/ekrishgupta/navtex/internal/latex"
)

func TestErrorModal_Navigation(t *testing.T) {
	em := NewErrorModal()

	entries := []latex.LogEntry{
		{Severity: "error", Line: 10, Message: "Error 1"},
		{Severity: "warning", Line: 15, Message: "Warning 1"},
		{Severity: "error", Line: 20, Message: "Error 2"},
	}

	em.Show(entries)

	if em.cursor != 0 {
		t.Errorf("Expected cursor to start at 0, got %d", em.cursor)
	}

	// Move down
	em.MoveDown()
	if em.cursor != 1 {
		t.Errorf("Expected cursor at 1 after MoveDown, got %d", em.cursor)
	}

	// Move up
	em.MoveUp()
	if em.cursor != 0 {
		t.Errorf("Expected cursor at 0 after MoveUp, got %d", em.cursor)
	}

	// Move up past bound
	em.MoveUp()
	if em.cursor != 0 {
		t.Errorf("Expected cursor to stay bounded at 0, got %d", em.cursor)
	}

	// Move down past bound
	em.MoveDown() // 1
	em.MoveDown() // 2
	em.MoveDown() // try 3
	if em.cursor != 2 {
		t.Errorf("Expected cursor to stay bounded at max index 2, got %d", em.cursor)
	}

	selected := em.SelectedEntry()
	if selected == nil || selected.Line != 20 {
		t.Errorf("Expected to select entry on line 20, got %v", selected)
	}
}
