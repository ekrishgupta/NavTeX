package ui

import (
	"github.com/ekrishgupta/navtex/internal/core"
)

// ── Messages ──

// ScannedMsg is sent when a directory scan completes.
type ScannedMsg struct{ Files *core.ProjectFiles }

// BuildFinishedMsg is sent when a LaTeX build completes.
type BuildFinishedMsg struct {
	Result *core.CompileResult
	Err    error
}

// LogParsedMsg is sent when a .log file has been parsed.
type LogParsedMsg struct {
	Entries []core.LogEntry
	Err     error
}

// CleanedMsg is sent when auxiliary files have been purged.
type CleanedMsg struct {
	Files []string
	Err   error
}

// ErrorMsg is sent for general asynchronous errors.
type ErrorMsg struct{ Err error }

// EditorClosedMsg is sent when the external editor process returns.
type EditorClosedMsg struct {
	Err error
}

// FileEventMsg is sent when the filesystem watcher detects a modification.
type FileEventMsg struct {
	Name string
}
