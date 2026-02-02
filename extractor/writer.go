package extractor

import (
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// SetTitle updates just the Title metadata.
func SetTitle(path string, title string) error {
	return setProperty(path, "Title", title)
}

// SetAuthor updates just the Author metadata.
func SetAuthor(path string, author string) error {
	return setProperty(path, "Author", author)
}

// SetNormalizedFlag sets the 'Normalized' metadata field to 'true'.
func SetNormalizedFlag(path string) error {
	return setProperty(path, "Normalized", "true")
}

// setProperty is a helper to update a specific key in the Info dictionary.
func setProperty(path string, key string, value string) error {
	// 1. Read context
	ctx, err := api.ReadContextFile(path)
	if err != nil {
		return err
	}

	// 2. Ensure Info dict exists
	if ctx.Info == nil {
		d := make(types.Dict)
		indRef, err := ctx.XRefTable.IndRefForNewObject(d)
		if err != nil {
			return err
		}
		ctx.Info = indRef
	}

	// 3. Dereference
	d, err := ctx.XRefTable.DereferenceDict(*ctx.Info)
	if err != nil {
		return err
	}

	// 4. Update Key
	d[key] = types.StringLiteral(value)

	// 5. Write back
	return api.WriteContextFile(ctx, path)
}
