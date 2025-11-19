package feedbin

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

// GetImports fetches all imports.
func (c *Client) GetImports(ctx context.Context) ([]Import, *Response, error) {
	req, err := c.NewRequest("GET", "imports.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var imports []Import
	resp, err := c.Do(ctx, req, &imports)
	if err != nil {
		return nil, resp, err
	}

	return imports, resp, nil
}

// GetImport fetches a single import by ID.
func (c *Client) GetImport(ctx context.Context, id int64) (*Import, *Response, error) {
	u := fmt.Sprintf("imports/%d.json", id)
	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var imp Import
	resp, err := c.Do(ctx, req, &imp)
	if err != nil {
		return nil, resp, err
	}

	return &imp, resp, nil
}

// CreateImport creates a new import from an OPML XML string.
func (c *Client) CreateImport(ctx context.Context, opmlContent []byte) (*Import, *Response, error) {
	req, err := c.NewRequest("POST", "imports.json", nil)
	if err != nil {
		return nil, nil, err
	}

	// Manually set body and content type for XML
	req.Body = io.NopCloser(bytes.NewReader(opmlContent))
	req.ContentLength = int64(len(opmlContent))
	req.Header.Set("Content-Type", "text/xml")

	var imp Import
	resp, err := c.Do(ctx, req, &imp)
	if err != nil {
		return nil, resp, err
	}

	return &imp, resp, nil
}
