package system

import (
	"fmt"
	"os"
	"os/exec"
)

// EditorCmd returns an exec.Cmd to open a file in the user's preferred editor.
// If lineNum is > 0, it attempts to open the editor at that specific line
// (supported for vi, vim, nvim, nano, code).
func EditorCmd(path string, lineNum int) (*exec.Cmd, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // Fallback
	}

	var args []string

	if lineNum > 0 {
		switch editor {
		case "vi", "vim", "nvim", "nano":
			args = append(args, fmt.Sprintf("+%d", lineNum))
		case "code":
			args = append(args, "--goto", fmt.Sprintf("%s:%d", path, lineNum))
			return exec.Command(editor, args...), nil // 'code --goto file:line'
		}
	}

	args = append(args, path)
	return exec.Command(editor, args...), nil
}
