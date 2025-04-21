package feedbin

import (
	"fmt"
	"net/http"
)

// SavedSearchesService handles saved search operations
type SavedSearchesService struct {
	client *Client
}

// List returns all saved searches
func (s *SavedSearchesService) List() ([]SavedSearch, error) {
	req, err := s.client.newRequest(http.MethodGet, "/saved_searches.json", nil)
	if err != nil {
		return nil, err
	}

	var searches []SavedSearch
	_, err = s.client.do(req, &searches)
	if err != nil {
		return nil, err
	}

	return searches, nil
}

// Get returns a specific saved search
func (s *SavedSearchesService) Get(id int) (*SavedSearch, error) {
	path := fmt.Sprintf("/saved_searches/%d.json", id)

	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	search := new(SavedSearch)
	_, err = s.client.do(req, search)
	if err != nil {
		return nil, err
	}

	return search, nil
}

// Create creates a new saved search
func (s *SavedSearchesService) Create(name, query string) (*SavedSearch, error) {
	createReq := SavedSearchCreateRequest{
		Name:  name,
		Query: query,
	}

	req, err := s.client.newRequest(http.MethodPost, "/saved_searches.json", createReq)
	if err != nil {
		return nil, err
	}

	search := new(SavedSearch)
	_, err = s.client.do(req, search)
	if err != nil {
		return nil, err
	}

	return search, nil
}

// Update updates a saved search
func (s *SavedSearchesService) Update(id int, name, query string) (*SavedSearch, error) {
	path := fmt.Sprintf("/saved_searches/%d.json", id)

	updateReq := SavedSearchUpdateRequest{
		Name:  name,
		Query: query,
	}

	req, err := s.client.newRequest(http.MethodPatch, path, updateReq)
	if err != nil {
		return nil, err
	}

	search := new(SavedSearch)
	_, err = s.client.do(req, search)
	if err != nil {
		return nil, err
	}

	return search, nil
}

// Delete deletes a saved search
func (s *SavedSearchesService) Delete(id int) error {
	path := fmt.Sprintf("/saved_searches/%d.json", id)

	req, err := s.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	return err
}
