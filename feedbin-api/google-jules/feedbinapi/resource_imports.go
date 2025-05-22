package feedbinapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// ImportsService handles operations related to OPML imports.
type ImportsService struct {
	client *Client
}

// ListImports retrieves all imports for the authenticated user.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/imports.md#get-v2importsjson
func (s *ImportsService) ListImports(ctx context.Context) ([]Import, *Response, error) {
	path := "imports.json"
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ListImports request: %w", err)
	}

	var imports []Import
	resp, err := s.client.do(req, &imports)
	if err != nil {
		return nil, resp, err
	}
	return imports, resp, nil
}

// GetImport retrieves a specific import by its ID, including the status of its items.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/imports.md#get-v2imports1json
func (s *ImportsService) GetImport(ctx context.Context, importID int64) (*Import, *Response, error) {
	if importID <= 0 {
		return nil, nil, fmt.Errorf("importID must be a positive integer")
	}
	path := fmt.Sprintf("imports/%d.json", importID)
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GetImport request: %w", err)
	}

	var imp Import
	resp, err := s.client.do(req, &imp)
	if err != nil {
		return nil, resp, err
	}
	return &imp, resp, nil
}

// CreateImport initiates a new import from an OPML file.
// The opmlData should be an io.Reader providing the XML content of the OPML file.
// The Content-Type header must be "text/xml".
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/imports.md#post-v2importsjson
func (s *ImportsService) CreateImport(ctx context.Context, opmlData io.Reader) (*Import, *Response, error) {
	if opmlData == nil {
		return nil, nil, fmt.Errorf("opmlData reader cannot be nil for CreateImport")
	}
	path := "imports.json"

	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil, opmlData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CreateImport request: %w", err)
	}
	// Override Content-Type for this specific request as per API docs
	req.Header.Set("Content-Type", "text/xml")

	var createdImport Import
	resp, err := s.client.do(req, &createdImport)
	if err != nil {
		return nil, resp, err
	}
	return &createdImport, resp, nil
}
