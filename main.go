package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pdf-title-extractor/extractor"
	"pdf-title-extractor/utils"
)

var (
	setTitleFlag  *string
	setAuthorFlag *string
	normalizeFlag *bool
)

func main() {
	setTitleFlag = flag.String("settitle", "", "Explicitly set the PDF title")
	setAuthorFlag = flag.String("setauthor", "", "Explicitly set the PDF author")
	normalizeFlag = flag.Bool("normalize", false, "Rename file based on Title and Author (Title-Author.pdf)")
	flag.Parse()

	files := getFiles()
	if len(files) == 0 {
		fmt.Println("Usage: pdf-title-extractor [options] <files...>")
		fmt.Println("       ls *.pdf | pdf-title-extractor [options]")
		flag.PrintDefaults()
		return
	}

	for _, file := range files {
		fmt.Printf("Processing: %s\n", file)
		if err := processFile(file); err != nil {
			fmt.Printf("Error processing '%s': %v. Moving to quarantine.\n", file, err)
			moveToError(file)
		}
		fmt.Println("---")
	}
}

func getFiles() []string {
	args := flag.Args()
	if len(args) > 0 {
		return args
	}

	// Check Stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		var files []string
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := strings.TrimSpace(scanner.Text())
			if text != "" {
				files = append(files, text)
			}
		}
		return files
	}

	return nil
}

func processFile(pdfPath string) (err error) {
	// 1. Panic Recovery
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("PANIC: %v", r)
		}
	}()

	// 2. Logic
	performedUpdate := false

	// MODE 1: Explicit Flags
	if *setTitleFlag != "" {
		fmt.Printf("Setting title to: '%s'...\n", *setTitleFlag)
		if err := extractor.SetTitle(pdfPath, *setTitleFlag); err != nil {
			return err
		}
		performedUpdate = true
	}

	if *setAuthorFlag != "" {
		fmt.Printf("Setting author to: '%s'...\n", *setAuthorFlag)
		if err := extractor.SetAuthor(pdfPath, *setAuthorFlag); err != nil {
			return err
		}
		performedUpdate = true
	}

	if performedUpdate {
		fmt.Println("Metadata updated successfully.")
		return nil
	}

	// MODE 3: Normalize (Rename)
	if *normalizeFlag {
		fmt.Println("Normalizing filename based on metadata...")
		meta, err := extractor.GetMetadata(pdfPath)
		if err != nil {
			return fmt.Errorf("reading metadata: %w", err)
		}

		if meta.Title == "" || meta.Author == "" {
			return fmt.Errorf("cannot normalize: missing Title or Author (Title='%s', Author='%s')", meta.Title, meta.Author)
		}

		cleanTitle := utils.Sanitize(meta.Title)
		cleanAuthor := utils.Sanitize(meta.Author)

		dir := filepath.Dir(pdfPath)
		ext := filepath.Ext(pdfPath)
		newName := fmt.Sprintf("%s-%s%s", cleanTitle, cleanAuthor, ext)
		desiredPath := filepath.Join(dir, newName)

		if desiredPath == pdfPath {
			fmt.Println("Filename is already normalized.")
			return nil
		}

		fmt.Println("Setting 'Normalized' flag...")
		if err := extractor.SetNormalizedFlag(pdfPath); err != nil {
			return fmt.Errorf("setting flag: %w", err)
		}

		finalPath := utils.GetUniquePath(desiredPath)
		fmt.Printf("Renaming '%s' to '%s'...\n", filepath.Base(pdfPath), filepath.Base(finalPath))

		if err := os.Rename(pdfPath, finalPath); err != nil {
			return fmt.Errorf("renaming: %w", err)
		}
		fmt.Println("File renamed successfully.")
		return nil
	}

	// MODE 2: Check & Auto-Set or View (Default)
	meta, err := extractor.GetMetadata(pdfPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning reading metadata for '%s': %v\n", pdfPath, err)
		// We can continue if we are just attempting auto-fix, or fail?
		// If read fails, auto-fix will likely fail too.
		return err
	}

	if meta.Title != "" {
		printMetadata(meta)
		return nil
	}

	// Auto-set
	fmt.Println("No metadata title found.")
	base := filepath.Base(pdfPath)
	ext := filepath.Ext(base)
	filenameTitle := base[0 : len(base)-len(ext)]
	filenameTitle = strings.ReplaceAll(filenameTitle, "_", " ")
	filenameTitle = strings.ReplaceAll(filenameTitle, "-", " ")

	fmt.Printf("Auto-setting title to: '%s'...\n", filenameTitle)
	if err := extractor.SetTitle(pdfPath, filenameTitle); err != nil {
		return err
	}

	defaultAuthor := "NIF"
	fmt.Printf("Auto-setting author to: '%s'...\n", defaultAuthor)
	if err := extractor.SetAuthor(pdfPath, defaultAuthor); err != nil {
		return err
	}

	fmt.Println("Metadata updated successfully.")
	return nil
}

func moveToError(path string) {
	dir := filepath.Dir(path)
	errorDir := filepath.Join(dir, "error")

	if err := os.MkdirAll(errorDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create error directory: %v\n", err)
		return
	}

	base := filepath.Base(path)
	dest := filepath.Join(errorDir, base)
	uniqueDest := utils.GetUniquePath(dest)

	fmt.Printf("Moving broken file to: %s\n", uniqueDest)
	if err := os.Rename(path, uniqueDest); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to move file to quarantine: %v\n", err)
	}
}

func printMetadata(meta extractor.PDFMetadata) {
	fmt.Println("--- Metadata Found ---")
	if meta.Title != "" {
		fmt.Printf("Title:             %s\n", meta.Title)
	}
	if meta.Author != "" {
		fmt.Printf("Author:            %s\n", meta.Author)
	}
	if meta.Subject != "" {
		fmt.Printf("Subject:           %s\n", meta.Subject)
	}
	if meta.Keywords != "" {
		fmt.Printf("Keywords:          %s\n", meta.Keywords)
	}
	if meta.CreationDate != "" {
		fmt.Printf("Creation Date:     %s\n", meta.CreationDate)
	}
	if meta.ModDate != "" {
		fmt.Printf("Modification Date: %s\n", meta.ModDate)
	}
	if meta.Normalized != "" {
		fmt.Printf("Normalized:        %s\n", meta.Normalized)
	}
	fmt.Println("----------------------")
}
