package core

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ProjectConfig holds settings defined in .navtex.yaml
type ProjectConfig struct {
	Engine     string   `yaml:"engine"`
	MasterFile string   `yaml:"master"`
	Ignores    []string `yaml:"ignores"`
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() ProjectConfig {
	return ProjectConfig{
		Engine:     "pdflatex",
		MasterFile: "",
	}
}

// LoadConfig attempts to read .navtex.yaml from the project root.
// Returns DefaultConfig if not found.
func LoadConfig(dir string) ProjectConfig {
	config := DefaultConfig()
	path := filepath.Join(dir, ".navtex.yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		return config // Return default silently if missing
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config // Assume default if broken
	}

	return config
}
