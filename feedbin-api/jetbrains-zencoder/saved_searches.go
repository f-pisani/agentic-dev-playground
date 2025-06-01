package feedbin

import (
	"fmt"
	"net/http"
)

// GetSavedSearches retrieves all saved searches for the authenticated user.
func (c *Client) GetSavedSearches() ([]SavedSearch, error) {
	// Make the request
	resp, err := c.doRequest(http.MethodGet, "/saved_searches.json", nil, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var searches []SavedSearch
	if err := parseResponse(resp, &searches); err != nil {
		return nil, err
	}

	return searches, nil
}

// GetSavedSearch retrieves a specific saved search by ID.
func (c *Client) GetSavedSearch(id int) (*SavedSearch, error) {
	// Make the request
	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("/saved_searches/%d.json", id), nil, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var search SavedSearch
	if err := parseResponse(resp, &search); err != nil {
		return nil, err
	}

	return &search, nil
}

// CreateSavedSearch creates a new saved search.
func (c *Client) CreateSavedSearch(name, query string) (*SavedSearch, error) {
	// Create the request body
	reqBody := SavedSearchRequest{
		Name:  name,
		Query: query,
	}

	// Make the request
	resp, err := c.doRequest(http.MethodPost, "/saved_searches.json", reqBody, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var search SavedSearch
	if err := parseResponse(resp, &search); err != nil {
		return nil, err
	}

	return &search, nil
}

// UpdateSavedSearch updates a saved search.
func (c *Client) UpdateSavedSearch(id int, name, query string) (*SavedSearch, error) {
	// Create the request body
	reqBody := SavedSearchRequest{
		Name:  name,
		Query: query,
	}

	// Make the request
	resp, err := c.doRequest(http.MethodPut, fmt.Sprintf("/saved_searches/%d.json", id), reqBody, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var search SavedSearch
	if err := parseResponse(resp, &search); err != nil {
		return nil, err
	}

	return &search, nil
}

// DeleteSavedSearch deletes a saved search by ID.
func (c *Client) DeleteSavedSearch(id int) error {
	// Make the request
	resp, err := c.doRequest(http.MethodDelete, fmt.Sprintf("/saved_searches/%d.json", id), nil, nil)
	if err != nil {
		return err
	}

	// Check for success (204 No Content)
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
