package core

import (
	"os"
	"strings"
	"testing"
)

func TestEditorCmd(t *testing.T) {
	// Temporarily override EDITOR
	originalEditor := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", originalEditor)

	cases := []struct {
		editor   string
		path     string
		line     int
		expected string // We will just check if strings are contained in args
	}{
		{"nvim", "main.tex", 42, "+42 main.tex"},
		{"code", "main.tex", 42, "--goto main.tex:42"},
		{"nano", "main.tex", 0, "main.tex"},  // No line number
		{"", "main.tex", 10, "+10 main.tex"}, // Fallback to vim
	}

	for _, tc := range cases {
		os.Setenv("EDITOR", tc.editor)
		cmd, _ := EditorCmd(tc.path, tc.line)

		argsJoined := strings.Join(cmd.Args[1:], " ")
		if !strings.Contains(argsJoined, tc.expected) {
			t.Errorf("For editor '%s', expected args to contain '%s', got '%s'", tc.editor, tc.expected, argsJoined)
		}
	}
}
