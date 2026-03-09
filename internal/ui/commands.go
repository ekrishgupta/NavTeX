package ui

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ekrishgupta/navtex/internal/latex"
	"github.com/ekrishgupta/navtex/internal/system"
)

// ── Commands ──

func (m Model) scanDirCmd(root string) tea.Cmd {
	return func() tea.Msg {
		pf, err := latex.ScanDirectory(root)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return ScannedMsg{Files: pf}
	}
}

func (m Model) compileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		res, err := m.compiler.Compile(path, m.rootPath, m.engine)
		return BuildFinishedMsg{Result: res, Err: err}
	}
}

func (m Model) parseLogCmd(path string) tea.Cmd {
	return func() tea.Msg {
		entries, err := latex.ParseLog(path)
		return LogParsedMsg{Entries: entries, Err: err}
	}
}

func (m Model) cleanCmd() tea.Cmd {
	return func() tea.Msg {
		removed, err := latex.Purge(m.rootPath)
		return CleanedMsg{Files: removed, Err: err}
	}
}

func (m Model) openPdfCmd() tea.Cmd {
	return func() tea.Msg {
		if m.projectFiles == nil || len(m.projectFiles.Output) == 0 {
			return nil
		}
		// Open the first output PDF
		latex.OpenPDF(m.projectFiles.Output[0].Path)
		return nil
	}
}

func (m Model) openEditorCmd(path string, lineNum int) tea.Cmd {
	c, err := system.EditorCmd(path, lineNum)
	if err != nil {
		return func() tea.Msg { return ErrorMsg{Err: err} }
	}
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return EditorClosedMsg{Err: err}
	})
}

func (m Model) listenForFileEventCmd() tea.Cmd {
	return func() tea.Msg {
		if m.watcher == nil {
			return nil
		}
		// Block until an event occurs
		event, ok := <-m.watcher.Events
		if !ok {
			return nil // watcher closed
		}
		return FileEventMsg{Name: event}
	}
}

func (m Model) runTexCountCmd(path string) tea.Cmd {
	return func() tea.Msg {
		if _, err := exec.LookPath("texcount"); err != nil {
			return TexCountFinishedMsg{Path: path, Err: err}
		}
		total, inText, inHeaders, inCaptions, err := latex.RunTexCount(path)
		return TexCountFinishedMsg{
			Path:       path,
			Total:      total,
			InText:     inText,
			InHeaders:  inHeaders,
			InCaptions: inCaptions,
			Err:        err,
		}
	}
}
