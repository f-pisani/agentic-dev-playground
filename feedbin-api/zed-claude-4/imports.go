package feedbin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// GetImports retrieves all imports for the authenticated user
func (c *Client) GetImports(ctx context.Context) ([]Import, error) {
	var imports []Import
	_, err := c.makeRequest(ctx, "GET", "/imports.json", nil, &imports)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

// GetImport retrieves a specific import by ID with status of all import items
func (c *Client) GetImport(ctx context.Context, id int) (*Import, error) {
	path := fmt.Sprintf("/imports/%d.json", id)

	var importObj Import
	_, err := c.makeRequest(ctx, "GET", path, nil, &importObj)
	if err != nil {
		return nil, err
	}

	return &importObj, nil
}

// CreateImport creates a new OPML import from XML content
func (c *Client) CreateImport(ctx context.Context, opmlContent string) (*Import, error) {
	if opmlContent == "" {
		return nil, &ValidationError{
			Field:   "opml_content",
			Message: "OPML content is required",
		}
	}

	req, err := c.newImportRequest(ctx, "POST", "/imports.json", opmlContent)
	if err != nil {
		return nil, err
	}

	var importObj Import
	_, err = c.do(req, &importObj)
	if err != nil {
		return nil, err
	}

	return &importObj, nil
}

// newImportRequest creates a new HTTP request for OPML import with proper headers
func (c *Client) newImportRequest(ctx context.Context, method, path string, opmlContent string) (*http.Request, error) {
	if err := c.credentials.Validate(); err != nil {
		return nil, err
	}

	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL path: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), strings.NewReader(opmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set required headers for OPML import
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "text/xml")

	// Set authentication
	setBasicAuth(req, c.credentials)

	return req, nil
}

// CreateImportFromFile creates a new OPML import from a file
func (c *Client) CreateImportFromFile(ctx context.Context, reader io.Reader) (*Import, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read OPML content: %v", err)
	}

	return c.CreateImport(ctx, string(content))
}

// Convenience methods without context (use background context)

func (c *Client) GetImportsWithoutContext() ([]Import, error) {
	return c.GetImports(context.Background())
}

func (c *Client) GetImportWithoutContext(id int) (*Import, error) {
	return c.GetImport(context.Background(), id)
}

func (c *Client) CreateImportWithoutContext(opmlContent string) (*Import, error) {
	return c.CreateImport(context.Background(), opmlContent)
}

func (c *Client) CreateImportFromFileWithoutContext(reader io.Reader) (*Import, error) {
	return c.CreateImportFromFile(context.Background(), reader)
}
