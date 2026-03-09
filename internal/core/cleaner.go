package core

import (
	"os"
	"path/filepath"
)

// PreviewPurge returns the list of auxiliary files that would be deleted.
func PreviewPurge(root string) ([]string, error) {
	pf, err := ScanDirectory(root)
	if err != nil {
		return nil, err
	}

	var targets []string
	for _, f := range pf.Auxiliary {
		rel, _ := filepath.Rel(pf.Root, f.Path)
		targets = append(targets, rel)
	}

	return targets, nil
}

// Purge deletes all auxiliary files from the given directory and returns
// the list of files that were removed.
func Purge(root string) ([]string, error) {
	pf, err := ScanDirectory(root)
	if err != nil {
		return nil, err
	}

	var removed []string
	for _, f := range pf.Auxiliary {
		if err := os.Remove(f.Path); err == nil {
			rel, _ := filepath.Rel(pf.Root, f.Path)
			removed = append(removed, rel)
		}
	}

	return removed, nil
}
