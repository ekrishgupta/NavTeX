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
