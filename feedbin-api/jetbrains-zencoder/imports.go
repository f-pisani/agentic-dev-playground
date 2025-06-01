package feedbin

import (
	"fmt"
	"net/http"
)

// GetImports retrieves all imports for the authenticated user.
func (c *Client) GetImports() ([]Import, error) {
	// Make the request
	resp, err := c.doRequest(http.MethodGet, "/imports.json", nil, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var imports []Import
	if err := parseResponse(resp, &imports); err != nil {
		return nil, err
	}

	return imports, nil
}

// GetImport retrieves a specific import by ID.
func (c *Client) GetImport(id int) (*Import, error) {
	// Make the request
	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("/imports/%d.json", id), nil, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var importObj Import
	if err := parseResponse(resp, &importObj); err != nil {
		return nil, err
	}

	return &importObj, nil
}

// CreateImport creates a new import from OPML data.
func (c *Client) CreateImport(opml string) (*Import, error) {
	// Create the request body
	reqBody := ImportRequest{
		OPML: opml,
	}

	// Make the request
	resp, err := c.doRequest(http.MethodPost, "/imports.json", reqBody, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var importObj Import
	if err := parseResponse(resp, &importObj); err != nil {
		return nil, err
	}

	return &importObj, nil
}

// DeleteImport deletes an import by ID.
func (c *Client) DeleteImport(id int) error {
	// Make the request
	resp, err := c.doRequest(http.MethodDelete, fmt.Sprintf("/imports/%d.json", id), nil, nil)
	if err != nil {
		return err
	}

	// Check for success (204 No Content)
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
