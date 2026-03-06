package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ekrishgupta/navtex/internal/core"
)

// ── Commands ──

func (m Model) scanDirCmd(root string) tea.Cmd {
	return func() tea.Msg {
		pf, err := core.ScanDirectory(root)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return ScannedMsg{Files: pf}
	}
}

func (m Model) compileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		res, err := m.compiler.Compile(path, m.engine)
		return BuildFinishedMsg{Result: res, Err: err}
	}
}

func (m Model) parseLogCmd(path string) tea.Cmd {
	return func() tea.Msg {
		entries, err := core.ParseLog(path)
		return LogParsedMsg{Entries: entries, Err: err}
	}
}

func (m Model) cleanCmd() tea.Cmd {
	return func() tea.Msg {
		removed, err := core.Purge(m.rootPath)
		return CleanedMsg{Files: removed, Err: err}
	}
}

func (m Model) openPdfCmd() tea.Cmd {
	return func() tea.Msg {
		if m.projectFiles == nil || len(m.projectFiles.Output) == 0 {
			return nil
		}
		// Open the first output PDF
		core.OpenPDF(m.projectFiles.Output[0].Path)
		return nil
	}
}

func (m Model) openEditorCmd(path string, lineNum int) tea.Cmd {
	c, err := core.EditorCmd(path, lineNum)
	if err != nil {
		return func() tea.Msg { return ErrorMsg{Err: err} }
	}
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return EditorClosedMsg{Err: err}
	})
}
