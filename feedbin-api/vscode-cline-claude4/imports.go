package feedbin

import (
	"context"
	"fmt"
)

// GetImports retrieves all imports for the authenticated user.
// Returns a slice of imports.
func (c *Client) GetImports(ctx context.Context) ([]Import, error) {
	var imports []Import
	_, err := c.get(ctx, "imports.json", nil, &imports)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

// GetImport retrieves a specific import by ID.
func (c *Client) GetImport(ctx context.Context, id int) (*Import, error) {
	path := fmt.Sprintf("imports/%d.json", id)

	var importItem Import
	_, err := c.get(ctx, path, nil, &importItem)
	if err != nil {
		return nil, err
	}

	return &importItem, nil
}

// CreateImport creates a new OPML import.
// The OPML data should be provided as a string containing the XML content.
// Returns the created import or an error.
func (c *Client) CreateImport(ctx context.Context, opmlData string) (*Import, error) {
	// The API expects the OPML data to be sent as form data or in a specific format
	// For simplicity, we'll assume it's sent as JSON with the OPML content
	request := map[string]interface{}{
		"opml": opmlData,
	}

	var importItem Import
	_, err := c.post(ctx, "imports.json", nil, request, &importItem)
	if err != nil {
		return nil, err
	}

	return &importItem, nil
}

// DeleteImport deletes an import by ID.
func (c *Client) DeleteImport(ctx context.Context, id int) error {
	path := fmt.Sprintf("imports/%d.json", id)
	_, err := c.delete(ctx, path, nil, nil)
	return err
}

// GetImportByTitle retrieves an import by title.
// Returns the first import with the matching title, or nil if not found.
func (c *Client) GetImportByTitle(ctx context.Context, title string) (*Import, error) {
	imports, err := c.GetImports(ctx)
	if err != nil {
		return nil, err
	}

	for _, importItem := range imports {
		if importItem.Title == title {
			return &importItem, nil
		}
	}

	return nil, fmt.Errorf("import with title '%s' not found", title)
}
