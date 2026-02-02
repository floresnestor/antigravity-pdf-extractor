# PDF Title Extractor

A robust Go command-line tool to extract metadata, normalize filenames, and manage PDF titles and authors.

## Features

- **Metadata Extraction**: View Title, Author, Creation Date, etc.
- **Normalization**: Renames files to `Title-Author.pdf` format.
  - Sanitizes filenames (e.g., `+` -> `Plus`, `ñ` -> `gn`).
  - Flags files as `Normalized: true` in metadata.
  - Avoids collisions (appends `_1`, `_2` if file exists).
- **Edit Metadata**: Explicitly set Title and Author via CLI.
- **Batch Processing**: Supports piping (`ls *.pdf | ...`) and multiple arguments.
- **Fault Tolerance**: Automatic quarantine of corrupt files to an `error/` folder.

## Usage

### Build
```bash
go build -o pdf-extractor.exe
```

### Examples

**View Metadata:**
```bash
./pdf-extractor.exe myfile.pdf
```

**Normalize (Rename & Sanitize):**
```bash
./pdf-extractor.exe --normalize myfile.pdf
```

**Batch Normalize:**
```bash
ls *.pdf | ./pdf-extractor.exe --normalize
```

**Set Metadata Explicitly:**
```bash
./pdf-extractor.exe --settitle "My Report" --setauthor "John Doe" myfile.pdf
```
