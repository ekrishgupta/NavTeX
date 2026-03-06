package core

import (
	"os"
	"testing"
)

func TestParseLog_Error(t *testing.T) {
	content := `
This is some log text.
! LaTeX Error: Environment document undefined.

See the LaTeX manual or LaTeX Companion for explanation.
Type  H <return>  for immediate help.
 ...                                              
                                                  
l.5 \begin{document}
                    
`
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	os.WriteFile(tmpFile.Name(), []byte(content), 0o644)

	entries, err := ParseLog(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseLog failed: %v", err)
	}

	foundError := false
	for _, e := range entries {
		if e.Severity == "error" && e.Line == 5 {
			foundError = true
			break
		}
	}

	if !foundError {
		t.Errorf("Expected LaTeX error on line 5 not found in %v", entries)
	}
}
