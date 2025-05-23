package feedbinapi

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	// "github.com/google/go-querystring/query" // Not needed for current endpoints
)

// ImportsService handles operations related to OPML imports.
type ImportsService struct {
	client *Client
}

// NewImportsService creates a new service for import operations.
func NewImportsService(client *Client) *ImportsService {
	return &ImportsService{client: client}
}

// ImportListOptions specifies optional parameters for listing imports.
// The API doc for "GET /v2/imports.json" does not specify any query parameters.
type ImportListOptions struct {
	// No options defined in spec for this endpoint.
}

// List retrieves all import jobs.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/imports.md#get-imports
func (s *ImportsService) List(opts *ImportListOptions) ([]Import, *http.Response, error) {
	path := "imports.json"
	// if opts != nil { ... } // No query options specified in docs

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var imports []Import
	resp, err := s.client.Do(req, &imports)
	if err != nil {
		return nil, resp, err
	}
	return imports, resp, nil
}

// Create schedules an OPML file for import.
// This requires a multipart/form-data POST request with the OPML file.
// Parameter `opml_file` should contain the file data.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/imports.md#create-import
func (s *ImportsService) Create(opmlFilePath string) (*Import, *http.Response, error) {
	if opmlFilePath == "" {
		return nil, nil, fmt.Errorf("opmlFilePath is required")
	}

	file, err := os.Open(opmlFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open OPML file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("opml_file", filepath.Base(opmlFilePath))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to copy file to form: %w", err)
	}
	err = writer.Close() // This finalizes the multipart body
	if err != nil {
		return nil, nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	path := "imports.json"
	req, err := s.client.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, nil, err
	}

	// Set the Content-Type header for multipart form data.
	// The boundary will be set automatically by NewRequest if body is a *bytes.Buffer
	// and Content-Type is not application/json.
	// However, for multipart, it's crucial to set it correctly including the boundary.
	// The http.NewRequest function does not automatically set Content-Type for multipart.
	// We need to set it from our multipart writer.
	req.Header.Set("Content-Type", writer.FormDataContentType())


	var createdImport Import
	resp, err := s.client.Do(req, &createdImport)
	if err != nil {
		return nil, resp, err
	}
	return &createdImport, resp, nil
}
