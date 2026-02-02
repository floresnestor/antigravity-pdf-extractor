# PDF Title Extractor

A robust Go command-line tool to extract metadata, normalize filenames based on content, and manage PDF titles and authors efficiently.

## Features

- **Metadata Extraction**: View Title, Author, Subject, Keywords, Creation Date, and the custom `Normalized` flag.
- **Normalization (Renaming)**:
  - Renames files to `SanitizedTitle-SanitizedAuthor.pdf` (e.g., `My_Doc-Me.pdf`).
  - **Sanitization Rules**:
    - Replaces `+` with `Plus` and `ñ` with `gn`.
    - Replaces spaces with underscores (`_`).
    - Removes all non-alphanumeric characters (keeps only A-Z, 0-9, `_`).
  - **Collision Avoidance**: If a file with the target name exists, appends a suffix (e.g., `_1`, `_2`) to prevent overwriting.
  - **Metadata Tagging**: Sets `Normalized: true` in the PDF metadata.
- **Edit Metadata**: Explicitly set Title and Author via CLI flags (`--settitle`, `--setauthor`).
- **Auto-Logic**: If a file has no title, the tool can automatically set the Title from the filename and Author to "NIF".
- **Batch Processing**: Supports processing multiple files via arguments or piping.
- **Fault Tolerance (Quarantine)**: Automatically moves corrupt or failing files to an `error/` directory, ensuring the batch process continues smoothly.

## Build and Run

### Prerequisites
- Go installed (1.16+ recommended).

### Build
```powershell
go build -o pdf-extractor.exe
```

## Usage

### 1. View Metadata
Simply run the tool on a PDF to view its metadata.
```powershell
./pdf-extractor.exe document.pdf
```

### 2. Normalize (Rename & Sanitize)
Renames the file based on its Title and Author metadata.
```powershell
./pdf-extractor.exe --normalize report.pdf
```
*Output: Renames `report.pdf` to `Title-Author.pdf`.*

### 3. Edit Metadata Explicitly
Set specific metadata fields.
```powershell
./pdf-extractor.exe --settitle "Annual Report" --setauthor "Jane Doe" report.pdf
```

### 4. Batch Processing
The tool supports batch operations using standard shell globbing or piping.

**Process multiple specific files:**
```powershell
./pdf-extractor.exe --normalize file1.pdf file2.pdf
```

**Process all PDFs in a folder (PowerShell/CMD):**
```powershell
ls *.pdf | ./pdf-extractor.exe --normalize
```

## Robustness & Error Handling
- **Quarantine**: If a file is corrupt (e.g., not a valid PDF), it is moved to an `./error/` subfolder.
- **Safe Renaming**: The tool never overwrites existing files during normalization. It generates unique names (e.g., `Report-Author_1.pdf`) instead.

## License
Open Source.
