// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"fmt"
	"net/http"
)

// SavedSearchService handles communication with the saved search related
// methods of the Feedbin API.
type SavedSearchService struct {
	client *Client
}

// GetSavedSearches retrieves all saved searches.
func (s *SavedSearchService) GetSavedSearches() ([]SavedSearch, error) {
	req, err := s.client.NewRequest("GET", "/saved_searches.json", nil)
	if err != nil {
		return nil, err
	}

	var searches []SavedSearch
	_, err = s.client.Do(req, &searches)
	if err != nil {
		return nil, err
	}

	return searches, nil
}

// GetSavedSearch retrieves a specific saved search.
func (s *SavedSearchService) GetSavedSearch(id int) (*SavedSearch, error) {
	path := fmt.Sprintf("/saved_searches/%d.json", id)
	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var search SavedSearch
	_, err = s.client.Do(req, &search)
	if err != nil {
		return nil, err
	}

	return &search, nil
}

// CreateSavedSearch creates a new saved search.
func (s *SavedSearchService) CreateSavedSearch(name, query string) (*SavedSearch, error) {
	body := map[string]interface{}{
		"name":  name,
		"query": query,
	}

	req, err := s.client.NewRequest("POST", "/saved_searches.json", body)
	if err != nil {
		return nil, err
	}

	var search SavedSearch
	_, err = s.client.Do(req, &search)
	if err != nil {
		return nil, err
	}

	return &search, nil
}

// UpdateSavedSearch updates an existing saved search.
func (s *SavedSearchService) UpdateSavedSearch(id int, name, query string) (*SavedSearch, error) {
	body := map[string]interface{}{
		"name":  name,
		"query": query,
	}

	path := fmt.Sprintf("/saved_searches/%d.json", id)
	req, err := s.client.NewRequest("PATCH", path, body)
	if err != nil {
		return nil, err
	}

	var search SavedSearch
	_, err = s.client.Do(req, &search)
	if err != nil {
		return nil, err
	}

	return &search, nil
}

// DeleteSavedSearch deletes a saved search.
func (s *SavedSearchService) DeleteSavedSearch(id int) error {
	path := fmt.Sprintf("/saved_searches/%d.json", id)
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