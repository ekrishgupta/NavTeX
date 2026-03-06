package core

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/main.tex.tmpl
var mainTexTemplate string

//go:embed templates/refs.bib.tmpl
var refsBibTemplate string

//go:embed templates/gitignore.tmpl
var gitignoreTemplate string

// CreateProject scaffolds a new LaTeX project in the given directory.
func CreateProject(root, title, author, docclass string) error {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolving path: %w", err)
	}

	if docclass == "" {
		docclass = "article"
	}

	// Create images directory
	imagesDir := filepath.Join(absRoot, "images")
	if err := os.MkdirAll(imagesDir, 0o755); err != nil {
		return fmt.Errorf("creating images directory: %w", err)
	}

	// Write main.tex
	mainContent := strings.ReplaceAll(mainTexTemplate, "{{TITLE}}", title)
	mainContent = strings.ReplaceAll(mainContent, "{{AUTHOR}}", author)
	mainContent = strings.ReplaceAll(mainContent, "{{DOCCLASS}}", docclass)
	mainContent = strings.ReplaceAll(mainContent, "{{DATE}}", "\\today")

	if err := writeIfNotExists(filepath.Join(absRoot, "main.tex"), mainContent); err != nil {
		return err
	}

	// Write refs.bib
	if err := writeIfNotExists(filepath.Join(absRoot, "refs.bib"), refsBibTemplate); err != nil {
		return err
	}

	// Write .gitignore
	if err := writeIfNotExists(filepath.Join(absRoot, ".gitignore"), gitignoreTemplate); err != nil {
		return err
	}

	// Create a .gitkeep in images/ so git tracks the empty directory
	gitkeep := filepath.Join(imagesDir, ".gitkeep")
	if err := writeIfNotExists(gitkeep, ""); err != nil {
		return err
	}

	return nil
}

// writeIfNotExists writes content to a file only if the file doesn't already exist.
func writeIfNotExists(path, content string) error {
	if _, err := os.Stat(path); err == nil {
		return nil // File already exists, skip
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", filepath.Base(path), err)
	}
	return nil
}
