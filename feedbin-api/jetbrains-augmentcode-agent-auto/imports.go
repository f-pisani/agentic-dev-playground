package feedbin

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ImportsService handles communication with the imports related
// methods of the Feedbin API.
type ImportsService struct {
	client *Client
}

// List returns all imports.
func (s *ImportsService) List() ([]*Import, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/v2/imports.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var imports []*Import
	resp, err := s.client.Do(req, &imports)
	if err != nil {
		return nil, resp, err
	}

	return imports, resp, nil
}

// Get returns a single import.
func (s *ImportsService) Get(id int64) (*Import, *http.Response, error) {
	u := fmt.Sprintf("/v2/imports/%d.json", id)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	imp := new(Import)
	resp, err := s.client.Do(req, imp)
	if err != nil {
		return nil, resp, err
	}

	return imp, resp, nil
}

// Create creates a new import from an OPML file.
func (s *ImportsService) Create(opmlContent string) (*Import, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/v2/imports.json", nil)
	if err != nil {
		return nil, nil, err
	}

	// Set the Content-Type header to text/xml
	req.Header.Set("Content-Type", "text/xml")

	// Set the request body to the OPML content
	req.Body = io.NopCloser(strings.NewReader(opmlContent))

	imp := new(Import)
	resp, err := s.client.Do(req, imp)
	if err != nil {
		return nil, resp, err
	}

	return imp, resp, nil
}
