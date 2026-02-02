package utils

import (
	"regexp"
	"strings"
)

// Sanitize creates a safe filename string from input.
// Rules:
// 1. + -> Plus
// 2. ñ -> gn, Ñ -> Gn
// 3. Space -> _
// 4. Remove all non-alphanumeric (English only) and non-underscore characters.
func Sanitize(input string) string {
	// 1. Replace + with Plus
	s := strings.ReplaceAll(input, "+", "Plus")

	// 2. Replace ñ/Ñ
	s = strings.ReplaceAll(s, "ñ", "gn")
	s = strings.ReplaceAll(s, "Ñ", "Gn")

	// 3. Replace spaces with _
	s = strings.ReplaceAll(s, " ", "_")

	// 4. Remove non-alphanumeric (A-Z, a-z, 0-9) and non-underscore
	// We use a regex to find invalid characters and replace them with empty string.
	// Valid: [a-zA-Z0-9_]
	re := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	s = re.ReplaceAllString(s, "")

	return s
}
