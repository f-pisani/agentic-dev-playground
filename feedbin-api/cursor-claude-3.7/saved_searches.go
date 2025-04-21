package feedbin

import (
	"fmt"
	"net/http"
)

// SavedSearchesService handles communication with the saved searches related
// endpoints of the Feedbin API
type SavedSearchesService struct {
	client *Client
}

// SavedSearchRequest represents a request to create or update a saved search
type SavedSearchRequest struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// List returns all saved searches
// https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#get-saved-searches
func (s *SavedSearchesService) List() ([]*SavedSearch, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "saved_searches.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var searches []*SavedSearch
	resp, err := s.client.Do(req, &searches)
	if err != nil {
		return nil, resp, err
	}

	return searches, resp, nil
}

// Get returns a specific saved search
// https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#get-saved-search
func (s *SavedSearchesService) Get(id int) (*SavedSearch, *http.Response, error) {
	url := fmt.Sprintf("saved_searches/%d.json", id)

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	search := new(SavedSearch)
	resp, err := s.client.Do(req, search)
	if err != nil {
		return nil, resp, err
	}

	return search, resp, nil
}

// Create creates a new saved search
// https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#create-saved-search
func (s *SavedSearchesService) Create(name, query string) (*SavedSearch, *http.Response, error) {
	request := &SavedSearchRequest{
		Name:  name,
		Query: query,
	}

	req, err := s.client.NewRequest(http.MethodPost, "saved_searches.json", request)
	if err != nil {
		return nil, nil, err
	}

	search := new(SavedSearch)
	resp, err := s.client.Do(req, search)
	if err != nil {
		return nil, resp, err
	}

	return search, resp, nil
}

// Update updates an existing saved search
// https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#update-saved-search
func (s *SavedSearchesService) Update(id int, name, query string) (*SavedSearch, *http.Response, error) {
	url := fmt.Sprintf("saved_searches/%d.json", id)

	request := &SavedSearchRequest{
		Name:  name,
		Query: query,
	}

	req, err := s.client.NewRequest(http.MethodPatch, url, request)
	if err != nil {
		return nil, nil, err
	}

	search := new(SavedSearch)
	resp, err := s.client.Do(req, search)
	if err != nil {
		return nil, resp, err
	}

	return search, resp, nil
}

// Delete removes a saved search
// https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#delete-saved-search
func (s *SavedSearchesService) Delete(id int) (*http.Response, error) {
	url := fmt.Sprintf("saved_searches/%d.json", id)

	req, err := s.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
