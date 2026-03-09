package latex

import (
	"os"
	"testing"
)

func TestScanDirectory_Empty(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "navtex-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	pf, err := ScanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	if pf.Total() != 0 {
		t.Errorf("Expected 0 files, got %d", pf.Total())
	}
}

func TestScanDirectory_Source(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "navtex-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	os.WriteFile(tmpDir+"/main.tex", []byte(""), 0o644)

	pf, err := ScanDirectory(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(pf.Source) != 1 {
		t.Errorf("Expected 1 source file, got %d", len(pf.Source))
	}
	if pf.Source[0].Name != "main.tex" {
		t.Errorf("Expected main.tex, got %s", pf.Source[0].Name)
	}
}

func TestScanDirectory_Assets(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "navtex-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	os.WriteFile(tmpDir+"/fig1.png", []byte(""), 0o644)
	os.WriteFile(tmpDir+"/photo.jpg", []byte(""), 0o644)

	pf, err := ScanDirectory(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(pf.Assets) != 2 {
		t.Errorf("Expected 2 assets, got %d", len(pf.Assets))
	}
}

func TestScanDirectory_Output(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "navtex-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	os.WriteFile(tmpDir+"/main.tex", []byte(""), 0o644)
	os.WriteFile(tmpDir+"/main.pdf", []byte(""), 0o644)
	os.WriteFile(tmpDir+"/other.pdf", []byte(""), 0o644)

	pf, err := ScanDirectory(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(pf.Output) != 1 {
		t.Fatalf("Expected 1 output file, got %d", len(pf.Output))
	}
	if pf.Output[0].Name != "main.pdf" {
		t.Errorf("Expected main.pdf, got %s", pf.Output[0].Name)
	}
	// other.pdf should be in Assets if no .tex matches (assuming pdf is asset if not output)
	foundOther := false
	for _, a := range pf.Assets {
		if a.Name == "other.pdf" {
			foundOther = true
			break
		}
	}
	if !foundOther {
		t.Errorf("Expected other.pdf to be classified as asset")
	}
}
