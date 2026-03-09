package latex

import (
	"os"
	"testing"
)

func TestTexMetadata_Basic(t *testing.T) {
	content := `
\documentclass[12pt]{article}
\usepackage{amsmath}
\usepackage{graphicx}
\title{Test Document}
\author{Author Name}
\begin{document}
Hello world.
\end{document}
`
	tmpFile, err := os.CreateTemp("", "test-*.tex")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	os.WriteFile(tmpFile.Name(), []byte(content), 0o644)

	meta, err := TexMetadata(tmpFile.Name())
	if err != nil {
		t.Fatalf("TexMetadata failed: %v", err)
	}

	if meta.Title != "Test Document" {
		t.Errorf("Expected title 'Test Document', got '%s'", meta.Title)
	}
	if meta.Author != "Author Name" {
		t.Errorf("Expected author 'Author Name', got '%s'", meta.Author)
	}
	if meta.DocumentClass != "article" {
		t.Errorf("Expected class 'article', got '%s'", meta.DocumentClass)
	}
}

func TestTexMetadata_WordCount(t *testing.T) {
	content := `
\documentclass{article}
\begin{document}
One two three four five.
\end{document}
`
	tmpFile, err := os.CreateTemp("", "test-*.tex")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	os.WriteFile(tmpFile.Name(), []byte(content), 0o644)

	meta, err := TexMetadata(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if meta.WordCount != 5 {
		t.Errorf("Expected 5 words, got %d", meta.WordCount)
	}
}

func TestBibMetadata_Basic(t *testing.T) {
	content := `
@article{key1,
  author = {Smith, John},
  title = {A Great Paper},
  year = {2023},
  journal = {Journal of Testing}
}
`
	tmpFile, err := os.CreateTemp("", "test-*.bib")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	os.WriteFile(tmpFile.Name(), []byte(content), 0o644)

	entries, err := BibMetadata(tmpFile.Name())
	if err != nil {
		t.Fatalf("BibMetadata failed: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}
	if entries[0].Key != "key1" {
		t.Errorf("Expected key 'key1', got '%s'", entries[0].Key)
	}
	if entries[0].Authors != "Smith, John" {
		t.Errorf("Expected author 'Smith, John', got '%s'", entries[0].Authors)
	}
}

func TestBibMetadata_RichFields(t *testing.T) {
	content := `
@article{key2,
  title = {Rich Paper},
  doi = {10.1234/test},
  url = {https://example.com},
  keywords = {latex, tui, go},
  abstract = {This is a very interesting abstract.}
}
`
	tmpFile, err := os.CreateTemp("", "test-*.bib")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	os.WriteFile(tmpFile.Name(), []byte(content), 0o644)

	entries, err := BibMetadata(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	e := entries[0]
	if e.DOI != "10.1234/test" {
		t.Errorf("Expected DOI 10.1234/test, got %s", e.DOI)
	}
	if e.URL != "https://example.com" {
		t.Errorf("Expected URL https://example.com, got %s", e.URL)
	}
	if len(e.Keywords) != 3 {
		t.Errorf("Expected 3 keywords, got %d", len(e.Keywords))
	}
	if e.Abstract != "This is a very interesting abstract." {
		t.Errorf("Expected abstract, got '%s'", e.Abstract)
	}
}
