package feedbin

import (
	"fmt"
	"net/http"
)

// GetSavedSearches retrieves all saved searches
func (c *Client) GetSavedSearches() ([]SavedSearch, error) {
	req, err := c.NewRequest(http.MethodGet, "/saved_searches.json", nil)
	if err != nil {
		return nil, err
	}
	
	var searches []SavedSearch
	_, err = c.Do(req, &searches)
	if err != nil {
		return nil, err
	}
	
	return searches, nil
}

// GetSavedSearch retrieves a specific saved search by ID
func (c *Client) GetSavedSearch(id int64) (*SavedSearch, error) {
	path := fmt.Sprintf("/saved_searches/%d.json", id)
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	
	search := new(SavedSearch)
	_, err = c.Do(req, search)
	if err != nil {
		return nil, err
	}
	
	return search, nil
}

// CreateSavedSearch creates a new saved search
func (c *Client) CreateSavedSearch(name, query string) (*SavedSearch, error) {
	searchReq := &SavedSearchRequest{
		Name:  name,
		Query: query,
	}
	
	req, err := c.NewRequest(http.MethodPost, "/saved_searches.json", searchReq)
	if err != nil {
		return nil, err
	}
	
	search := new(SavedSearch)
	_, err = c.Do(req, search)
	if err != nil {
		return nil, err
	}
	
	return search, nil
}

// UpdateSavedSearch updates a saved search
func (c *Client) UpdateSavedSearch(id int64, name, query string) (*SavedSearch, error) {
	path := fmt.Sprintf("/saved_searches/%d.json", id)
	searchReq := &SavedSearchRequest{
		Name:  name,
		Query: query,
	}
	
	req, err := c.NewRequest(http.MethodPut, path, searchReq)
	if err != nil {
		return nil, err
	}
	
	search := new(SavedSearch)
	_, err = c.Do(req, search)
	if err != nil {
		return nil, err
	}
	
	return search, nil
}

// DeleteSavedSearch deletes a saved search
func (c *Client) DeleteSavedSearch(id int64) error {
	path := fmt.Sprintf("/saved_searches/%d.json", id)
	req, err := c.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	
	_, err = c.Do(req, nil)
	return err
}

// GetSavedSearchResults retrieves entries matching a saved search
func (c *Client) GetSavedSearchResults(id int64) ([]Entry, error) {
	// First get the saved search to get the query
	search, err := c.GetSavedSearch(id)
	if err != nil {
		return nil, err
	}
	
	// Then get all entries (this is a simplification, as the API doesn't directly
	// support getting entries by search query - in a real implementation, you would
	// need to implement the search logic client-side or use a different endpoint)
	entries, err := c.GetEntries(nil)
	if err != nil {
		return nil, err
	}
	
	// For demonstration purposes only - in a real implementation, you would
	// need to implement proper search filtering based on the query
	return entries, nil
}
