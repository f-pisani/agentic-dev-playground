// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"fmt"
	"net/http"
)

// ImportService handles communication with the import related
// methods of the Feedbin API.
type ImportService struct {
	client *Client
}

// GetImports retrieves all imports.
func (s *ImportService) GetImports() ([]Import, error) {
	req, err := s.client.NewRequest("GET", "/imports.json", nil)
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

// GetImport retrieves a specific import.
func (s *ImportService) GetImport(id int) (*Import, error) {
	path := fmt.Sprintf("/imports/%d.json", id)
	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var imp Import
	_, err = s.client.Do(req, &imp)
	if err != nil {
		return nil, err
	}

	return &imp, nil
}

// CreateImportFromURL creates a new import from a URL.
func (s *ImportService) CreateImportFromURL(opmlURL string) (*Import, error) {
	body := map[string]interface{}{
		"opml_url": opmlURL,
	}

	req, err := s.client.NewRequest("POST", "/imports.json", body)
	if err != nil {
		return nil, err
	}

	var imp Import
	_, err = s.client.Do(req, &imp)
	if err != nil {
		return nil, err
	}

	return &imp, nil
}

// CreateImportFromFile creates a new import from a file.
// Note: This is a placeholder as file uploads require multipart form data,
// which is not implemented in this client.
func (s *ImportService) CreateImportFromFile(opmlFile string) (*Import, error) {
	// This would require multipart form data handling
	return nil, fmt.Errorf("file upload not implemented")
}

// DeleteImport deletes an import.
func (s *ImportService) DeleteImport(id int) error {
	path := fmt.Sprintf("/imports/%d.json", id)
	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}