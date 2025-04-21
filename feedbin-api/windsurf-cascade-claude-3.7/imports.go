package feedbin

import (
	"fmt"
	"net/http"
	"strings"
)

// ImportsService handles operations related to OPML imports
type ImportsService struct {
	client *Client
}

// List returns all imports for the user
func (s *ImportsService) List() ([]Import, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/imports.json", nil)
	if err != nil {
		return nil, err
	}

	var imports []Import
	_, err = s.client.Do(req, &imports)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

// Get returns a single import by ID
func (s *ImportsService) Get(id int) (*Import, error) {
	path := fmt.Sprintf("/imports/%d.json", id)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	imp := new(Import)
	_, err = s.client.Do(req, imp)
	if err != nil {
		return nil, err
	}

	return imp, nil
}

// Create creates a new import from an OPML XML string
func (s *ImportsService) Create(opmlContent string) (*Import, error) {
	// Create a custom request since this endpoint requires XML content
	u, err := s.client.BaseURL.Parse("/imports.json")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(opmlContent))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(s.client.Username, s.client.Password)
	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", UserAgent)

	imp := new(Import)
	_, err = s.client.Do(req, imp)
	if err != nil {
		return nil, err
	}

	return imp, nil
}
