package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ekrishgupta/navtex/internal/latex"
)

// Inspector is the right-pane metadata viewer.
type Inspector struct {
	path     string
	category latex.FileCategory
	width    int
	height   int
	focused  bool
	style    InspectorBaseStyle

	// Bib selection
	selectedBibIdx int

	// Cached metadata
	texMeta      *latex.TexMeta
	bibMeta      []latex.BibEntry
	imageMeta    *latex.ImageMeta
	fileSize     int64
	err          error
	needsRefresh bool
}

// NewInspector creates a new inspector.
func NewInspector(style InspectorBaseStyle) Inspector {
	return Inspector{style: style}
}

// SetStyle switches the inspector's visual style (focused / blurred).
func (ins *Inspector) SetStyle(s InspectorBaseStyle) {
	ins.style = s
}

// SetSize sets the inspector dimensions.
func (ins *Inspector) SetSize(w, h int) {
	ins.width = w
	ins.height = h
}

// SetFocused sets focus state.
func (ins *Inspector) SetFocused(f bool) {
	ins.focused = f
}

// Refresh marks the inspector as needing a metadata reload.
func (ins *Inspector) Refresh() {
	ins.needsRefresh = true
}

// SetFile updates the inspector to show metadata for the given file.
func (ins *Inspector) SetFile(path string, cat latex.FileCategory) {
	if path == ins.path && !ins.needsRefresh {
		return // No change
	}

	ins.path = path
	ins.category = cat
	ins.needsRefresh = false
	ins.texMeta = nil
	ins.bibMeta = nil
	ins.imageMeta = nil
	ins.err = nil

	if path == "" {
		return
	}

	switch ins.category {
	case latex.CategorySource:
		ins.texMeta, ins.err = latex.TexMetadata(path)

	case latex.CategoryData:
		if strings.HasSuffix(strings.ToLower(path), ".bib") {
			ins.bibMeta, ins.err = latex.BibMetadata(path)
		}

	case latex.CategoryAssets:
		ins.imageMeta, ins.err = latex.ImageMetadata(path)

	default:
		// For other files, just show basic info
	}
}

// MoveBibUp moves the bibliography selection up.
func (ins *Inspector) MoveBibUp() {
	if len(ins.bibMeta) == 0 {
		return
	}
	ins.selectedBibIdx--
	if ins.selectedBibIdx < 0 {
		ins.selectedBibIdx = len(ins.bibMeta) - 1
	}
}

// MoveBibDown moves the bibliography selection down.
func (ins *Inspector) MoveBibDown() {
	if len(ins.bibMeta) == 0 {
		return
	}
	ins.selectedBibIdx++
	if ins.selectedBibIdx >= len(ins.bibMeta) {
		ins.selectedBibIdx = 0
	}
}

// SelectedBibKey returns the citekey of the currently selected bibliography entry.
func (ins Inspector) SelectedBibKey() string {
	if len(ins.bibMeta) == 0 || ins.selectedBibIdx < 0 || ins.selectedBibIdx >= len(ins.bibMeta) {
		return ""
	}
	return ins.bibMeta[ins.selectedBibIdx].Key
}

// View renders the inspector.
func (ins Inspector) View() string {
	s := ins.style
	innerW := ins.width - 2

	// Title bar
	titleBar := s.TitleBar.Width(innerW).Render("🔍 Inspector")

	var content string

	if ins.path == "" {
		content = lipgloss.Place(innerW, ins.height-3, lipgloss.Center, lipgloss.Center,
			DimText.Render("Select a file to inspect"),
		)
	} else if ins.err != nil {
		content = lipgloss.JoinVertical(lipgloss.Left,
			s.SectionTitle.Render(filepath.Base(ins.path)),
			"",
			s.ErrorText.Render("Error: "+ins.err.Error()),
		)
	} else {
		switch ins.category {
		case latex.CategorySource:
			content = ins.renderTexMeta()
		case latex.CategoryData:
			if ins.bibMeta != nil {
				content = ins.renderBibMeta()
			} else {
				content = ins.renderGeneric()
			}
		case latex.CategoryAssets:
			content = ins.renderImageMeta()
		default:
			content = ins.renderGeneric()
		}
	}

	return lipgloss.NewStyle().Width(ins.width).Height(ins.height).Margin(0, 1).
		Render(lipgloss.JoinVertical(lipgloss.Left, titleBar, content))
}

// renderTexMeta renders .tex file metadata.
func (ins Inspector) renderTexMeta() string {
	s := ins.style
	m := ins.texMeta
	lines := []string{
		s.SectionTitle.Render("📄 " + filepath.Base(ins.path)),
		"",
	}

	// Document info
	if m.Title != "" {
		lines = append(lines, s.MetaLabel.Render("Title")+s.MetaValue.Render(m.Title))
	}
	if m.Author != "" {
		lines = append(lines, s.MetaLabel.Render("Author")+s.MetaValue.Render(m.Author))
	}
	lines = append(lines, s.MetaLabel.Render("Class")+s.MetaValue.Render(m.DocumentClass))
	if m.ClassOptions != "" {
		lines = append(lines, s.MetaLabel.Render("Options")+s.MetaValue.Render(m.ClassOptions))
	}

	lines = append(lines,
		"",
		SeparatorLine(ins.width-8),
		"",
	)

	lines = append(lines, s.MetaLabel.Render("Word Count")+s.MetaValue.Render(fmt.Sprintf("%d", m.WordCount)))
	if m.WordsInText > 0 || m.WordsInHeaders > 0 || m.WordsInCaptions > 0 {
		lines = append(lines, "  "+s.MetaLabel.Render("Text")+s.MetaValue.Render(fmt.Sprintf("%d", m.WordsInText)))
		lines = append(lines, "  "+s.MetaLabel.Render("Headers")+s.MetaValue.Render(fmt.Sprintf("%d", m.WordsInHeaders)))
		lines = append(lines, "  "+s.MetaLabel.Render("Captions")+s.MetaValue.Render(fmt.Sprintf("%d", m.WordsInCaptions)))
	}

	// Packages
	if len(m.Packages) > 0 {
		lines = append(lines,
			"",
			SeparatorLine(ins.width-8),
			"",
			s.SectionTitle.Render(fmt.Sprintf("Packages (%d)", len(m.Packages))),
		)
		var pkgLine strings.Builder
		for i, pkg := range m.Packages {
			if i > 0 {
				pkgLine.WriteString(" ")
			}
			pkgLine.WriteString(s.PackageTag.Render(pkg.Name))
			if pkgLine.Len() > ins.width-8 {
				lines = append(lines, "   "+pkgLine.String())
				pkgLine.Reset()
			}
		}
		if pkgLine.Len() > 0 {
			lines = append(lines, "   "+pkgLine.String())
		}
	}

	return strings.Join(lines, "\n")
}

// renderBibMeta renders .bib file metadata in a tabular format.
func (ins Inspector) renderBibMeta() string {
	s := ins.style
	lines := []string{
		s.SectionTitle.Render("📚 " + filepath.Base(ins.path)),
		s.MetaValue.Render(fmt.Sprintf("   %d entries", len(ins.bibMeta))),
		"",
	}

	if len(ins.bibMeta) == 0 {
		lines = append(lines, DimText.Render("  No entries found"))
		return strings.Join(lines, "\n")
	}

	// Column widths
	maxAuth := 18
	maxTitle := ins.width - maxAuth - 14
	if maxTitle < 20 {
		maxTitle = 20
	}

	// Header
	header := fmt.Sprintf("  %-*s %-*s %-4s %-8s",
		maxAuth, "Authors", maxTitle, "Title", "Year", "Type")
	headerStyle := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	lines = append(lines, headerStyle.Render(header))
	lines = append(lines, DimText.Render("  "+SeparatorLine(min(ins.width-6, len(header)))))

	// Entries
	for i, entry := range ins.bibMeta {
		authors := truncate(entry.Authors, maxAuth)
		title := truncate(entry.Title, maxTitle)
		row := fmt.Sprintf("  %-*s %-*s %-4s %-8s",
			maxAuth, authors, maxTitle, title, entry.Year, entry.Type)

		if i == ins.selectedBibIdx && ins.focused {
			lines = append(lines, s.SelectedRow.Width(ins.width-4).Render("▸"+row[1:]))
		} else {
			lines = append(lines, s.UnselectedRow.Render(row))
		}

		if entry.DOI != "" {
			lines = append(lines, DimText.Render(fmt.Sprintf("    DOI: %s", entry.DOI)))
		}

		if len(entry.Keywords) > 0 {
			var kwLine strings.Builder
			kwLine.WriteString("    ")
			for j, kw := range entry.Keywords {
				if j > 0 {
					kwLine.WriteString(" ")
				}
				kwLine.WriteString(s.KeywordTag.Render(kw))
			}
			lines = append(lines, kwLine.String())
		}
	}

	return strings.Join(lines, "\n")
}

// renderImageMeta renders image file metadata.
func (ins Inspector) renderImageMeta() string {
	s := ins.style
	m := ins.imageMeta
	lines := []string{
		s.SectionTitle.Render("🖼  " + filepath.Base(ins.path)),
		"",
		s.MetaLabel.Render("Format") + s.MetaValue.Render(strings.ToUpper(m.Format)),
	}

	if m.Width > 0 && m.Height > 0 {
		lines = append(lines, s.MetaLabel.Render("Dimensions")+s.MetaValue.Render(fmt.Sprintf("%d × %d px", m.Width, m.Height)))
	}

	lines = append(lines, s.MetaLabel.Render("File Size")+s.MetaValue.Render(latex.FormatSize(m.Size)))

	return strings.Join(lines, "\n")
}

// renderGeneric renders basic file information.
func (ins Inspector) renderGeneric() string {
	s := ins.style
	name := filepath.Base(ins.path)
	ext := filepath.Ext(ins.path)

	return strings.Join([]string{
		s.SectionTitle.Render("📎 " + name),
		"",
		s.MetaLabel.Render("Extension") + s.MetaValue.Render(ext),
		s.MetaLabel.Render("Category") + s.MetaValue.Render(categoryName(ins.category)),
	}, "\n")
}

func categoryName(c latex.FileCategory) string {
	switch c {
	case latex.CategorySource:
		return "Source"
	case latex.CategoryData:
		return "Data/Bib"
	case latex.CategoryAssets:
		return "Asset"
	case latex.CategoryAuxiliary:
		return "Auxiliary"
	case latex.CategoryOutput:
		return "Output"
	default:
		return "Unknown"
	}
}
