package feedbin

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ImportsService handles communication with the import related
// methods of the Feedbin API.
type ImportsService service

// Create initiates a new OPML import.
// The body parameter should be an io.Reader containing the OPML XML data.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/imports.md#post-v2importsjson
func (s *ImportsService) Create(body io.Reader) (*Import, *http.Response, error) {
	// Need to create a request specifically setting Content-Type to text/xml
	rel, err := url.Parse("imports.json")
	if err != nil {
		return nil, nil, err
	}
	u := s.client.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(http.MethodPost, u.String(), body)
	if err != nil {
		return nil, nil, err
	}

	req.SetBasicAuth(s.client.Username, s.client.password)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "text/xml") // Crucial for OPML import
	req.Header.Set("Accept", "application/json")

	var imp Import
	resp, err := s.client.do(req, &imp)
	if err != nil {
		return nil, resp, err
	}

	return &imp, resp, nil
}

// List retrieves all import tasks for the user.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/imports.md#get-v2importsjson
func (s *ImportsService) List() ([]*Import, *http.Response, error) {
	req, err := s.client.newRequest(http.MethodGet, "imports.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var imports []*Import
	resp, err := s.client.do(req, &imports)
	if err != nil {
		return nil, resp, err
	}

	return imports, resp, nil
}

// Get retrieves the status of a specific import task by its ID.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/imports.md#get-v2imports1json
func (s *ImportsService) Get(importID int64) (*Import, *http.Response, error) {
	path := fmt.Sprintf("imports/%d.json", importID)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var imp Import
	resp, err := s.client.do(req, &imp)
	if err != nil {
		return nil, resp, err
	}

	return &imp, resp, nil
}
