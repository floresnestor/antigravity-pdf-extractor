package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetUniquePath checks if a file exists at the given path.
// If it does, it appends a number (e.g., _1, _2) to the filename until a unique path is found.
func GetUniquePath(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}

	dir := filepath.Dir(path)
	ext := filepath.Ext(path)
	name := strings.TrimSuffix(filepath.Base(path), ext)

	counter := 1
	for {
		newPath := filepath.Join(dir, fmt.Sprintf("%s_%d%s", name, counter, ext))
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
		counter++
	}
}
