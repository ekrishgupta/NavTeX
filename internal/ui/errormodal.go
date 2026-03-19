package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ekrishgupta/navtex/internal/latex"
)

// ErrorModal displays parsed build errors.
type ErrorModal struct {
	entries []latex.LogEntry
	visible bool
	width   int
	height  int
	scroll  int
	cursor  int
	modalW  int
	modalH  int

	// Cache
	headerStr string
}

// NewErrorModal creates a new error modal.
func NewErrorModal() ErrorModal {
	return ErrorModal{}
}

// Show displays the modal with the given log entries.
func (em *ErrorModal) Show(entries []latex.LogEntry) {
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
		if em.cursor < em.scroll {
			em.scroll = em.cursor
		}
	}
}

// MoveDown moves the cursor down.
func (em *ErrorModal) MoveDown() {
	if em.cursor < len(em.entries)-1 {
		em.cursor++
	}
}

// SelectedEntry returns the currently selected log entry.
func (em *ErrorModal) SelectedEntry() *latex.LogEntry {
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

	errors := latex.ErrorCount(em.entries)
	warnings := latex.WarningCount(em.entries)

	errorStyle := lipgloss.NewStyle().Foreground(ColorBrightRed).Bold(true)
	warningStyle := lipgloss.NewStyle().Foreground(ColorYellow)

	if modalW != em.modalW || modalH != em.modalH {
		em.modalW = modalW
		em.modalH = modalH

		lineCol := 6
		sevCol := 8
		header := fmt.Sprintf("  %-*s %-*s %s", lineCol, "Line", sevCol, "Severity", "Message")
		headerStyle := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
		em.headerStr = headerStyle.Render(header) + "\n" + DimText.Render("  "+SeparatorLine(modalW-8))
	}

	lineCol := 6
	sevCol := 8
	msgCol := modalW - lineCol - sevCol - 12

	title := ModalTitleBar.Render(fmt.Sprintf("Build Log — %d errors, %d warnings", errors, warnings))

	// Use red-accented border if errors present
	frame := ModalFrame
	if errors > 0 {
		frame = frame.BorderForeground(ColorBrightRed)
	}

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

	selectedStyle := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)

	for i := em.scroll; i < len(em.entries) && len(rows) < visibleRows; i++ {
		e := em.entries[i]
		lineStr := "—"
		if e.Line > 0 {
			lineStr = fmt.Sprintf("%d", e.Line)
		}

		msg := truncate(e.Message, msgCol)

		var sev string
		if e.Severity == "error" {
			sev = errorStyle.Render("✗ error")
		} else {
			sev = warningStyle.Render("⚠ warn ")
		}

		row := fmt.Sprintf("  %-*s %-*s %s", lineCol, lineStr, sevCol, sev, msg)

		if i == em.cursor {
			row = selectedStyle.Width(modalW - 4).Render("▸" + row[1:])
		} else {
			row = " " + row[1:]
		}

		rows = append(rows, row)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		em.headerStr,
		strings.Join(rows, "\n"),
		"",
		ModalHint.Render("Enter: jump to line │ Esc: close │ ↑/↓: move"),
	)

	modal := frame.Width(modalW).Render(content)
	return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, modal)
}
