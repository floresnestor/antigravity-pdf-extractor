package extractor

import (
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// PDFMetadata holds extracted metadata fields
type PDFMetadata struct {
	Title        string
	Author       string
	Subject      string
	Keywords     string
	CreationDate string
	ModDate      string
	Normalized   string
}

// GetMetadata extracts standard metadata from the PDF.
func GetMetadata(path string) (PDFMetadata, error) {
	var meta PDFMetadata

	// We only need the context to access the XRefTable and Info dictionary
	ctx, err := api.ReadContextFile(path)
	if err != nil {
		return meta, err
	}

	if ctx.Info == nil {
		return meta, nil // No metadata info
	}

	// ctx.Info is an IndirectRef, we must dereference it to get the Dict
	d, err := ctx.XRefTable.DereferenceDict(*ctx.Info)
	if err != nil {
		return meta, err
	}

	meta.Title = extractString(d, "Title")
	meta.Author = extractString(d, "Author")
	meta.Subject = extractString(d, "Subject")
	meta.Keywords = extractString(d, "Keywords")
	meta.CreationDate = extractString(d, "CreationDate")
	meta.ModDate = extractString(d, "ModDate")
	meta.Normalized = extractString(d, "Normalized")

	return meta, nil
}

func extractString(d types.Dict, key string) string {
	val, ok := d[key]
	if !ok {
		return ""
	}

	if s, ok := val.(types.StringLiteral); ok {
		str, _ := types.StringLiteralToString(s)
		return str
	}
	if s, ok := val.(types.HexLiteral); ok {
		str, _ := types.HexLiteralToString(s)
		return str
	}
	return ""
}
