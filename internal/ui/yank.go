package ui

import (
	"github.com/atotto/clipboard"
)

// YankToClipboard copies the given text to the system clipboard.
func YankToClipboard(text string) error {
	return clipboard.WriteAll(text)
}
