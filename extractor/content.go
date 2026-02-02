package extractor

import (
	"strings"

	"github.com/ledongthuc/pdf"
)

// SuggestTitleFromContent attempts to find a title from the text content of the first page.
func SuggestTitleFromContent(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if r.NumPage() == 0 {
		return "", nil
	}

	// Read first page only
	p := r.Page(1)
	if p.V.IsNull() {
		return "", nil
	}

	// Extract text
	content, err := p.GetPlainText(nil)
	if err != nil {
		return "", err
	}

	// Simple heuristic: First non-empty line
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 0 {
			// You could add heuristics here: e.g. "if line is too short, skip", or "if it looks like a page number"
			// For now, return the first significant text.
			return trimmed, nil
		}
	}

	return "", nil
}
