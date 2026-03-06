package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ekrishgupta/navtex/internal/core"
)

// ErrorModal displays parsed build errors.
type ErrorModal struct {
	entries []core.LogEntry
	visible bool
	width   int
	height  int
	scroll  int
	cursor  int
}

// NewErrorModal creates a new error modal.
func NewErrorModal() ErrorModal {
	return ErrorModal{}
}

// Show displays the modal with the given log entries.
func (em *ErrorModal) Show(entries []core.LogEntry) {
	em.entries = entries
	em.visible = true
	em.scroll = 0
	em.cursor = 0
}

// Hide closes the modal.
func (em *ErrorModal) Hide() {
	em.visible = false
}

// IsVisible returns whether the modal is shown.
func (em *ErrorModal) IsVisible() bool {
	return em.visible
}

// MoveUp moves the cursor up.
func (em *ErrorModal) MoveUp() {
	if em.cursor > 0 {
		em.cursor--
		// Adjust scroll if cursor moves above viewport
		if em.cursor < em.scroll {
			em.scroll = em.cursor
		}
	}
}

// MoveDown moves the cursor down.
func (em *ErrorModal) MoveDown() {
	if em.cursor < len(em.entries)-1 {
		em.cursor++
		// Adjust scroll. Let's estimate visible rows (modalH-6). We'll handle this in View, but roughly let's do 10 for now, or just calculate it cleanly later.
		// A simple way is to just let the View method correct the scroll if needed, but doing it here needs height. We'll add dynamic calculation.
	}
}

// SelectedEntry returns the currently selected log entry.
func (em *ErrorModal) SelectedEntry() *core.LogEntry {
	if em.cursor >= 0 && em.cursor < len(em.entries) {
		return &em.entries[em.cursor]
	}
	return nil
}

// View renders the error modal.
func (em ErrorModal) View(termWidth, termHeight int) string {
	if !em.visible {
		return ""
	}

	modalW := termWidth * 3 / 4
	modalH := termHeight * 3 / 4
	if modalW < 60 {
		modalW = 60
	}
	if modalH < 10 {
		modalH = 10
	}

	errors := core.ErrorCount(em.entries)
	warnings := core.WarningCount(em.entries)

	title := ModalTitle.Render(fmt.Sprintf("Build Log — %d errors, %d warnings", errors, warnings))

	// Header row
	lineCol := 6
	sevCol := 8
	msgCol := modalW - lineCol - sevCol - 12

	header := fmt.Sprintf("  %-*s %-*s %s", lineCol, "Line", sevCol, "Severity", "Message")
	headerLine := BibTableHeader.Render(header)
	separator := FileItemDim.Render("  " + strings.Repeat("─", modalW-8))

	var rows []string

	// Auto-adjust scroll to keep cursor in view
	visibleRows := modalH - 6
	if visibleRows < 1 {
		visibleRows = 1
	}
	if em.cursor < em.scroll {
		em.scroll = em.cursor
	} else if em.cursor >= em.scroll+visibleRows {
		em.scroll = em.cursor - visibleRows + 1
	}

	for i := em.scroll; i < len(em.entries) && len(rows) < visibleRows; i++ {
		e := em.entries[i]
		lineStr := "—"
		if e.Line > 0 {
			lineStr = fmt.Sprintf("%d", e.Line)
		}

		msg := truncate(e.Message, msgCol)

		var sev string
		if e.Severity == "error" {
			sev = ErrorText.Render("error")
		} else {
			sev = WarningText.Render("warning")
		}

		row := fmt.Sprintf("  %-*s %-*s %s", lineCol, lineStr, sevCol, sev, msg)

		// Highlight cursor
		if i == em.cursor {
			row = FileItemSelected.Width(modalW - 4).Render("▸" + row[1:])
		} else {
			row = " " + row[1:]
		}

		rows = append(rows, row)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		headerLine,
		separator,
		strings.Join(rows, "\n"),
		"",
		FileItemDim.Render("  Enter: jump to line │ Esc: close │ ↑/↓: move"),
	)

	modal := ModalBox.Width(modalW).Render(content)
	return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, modal)
}
