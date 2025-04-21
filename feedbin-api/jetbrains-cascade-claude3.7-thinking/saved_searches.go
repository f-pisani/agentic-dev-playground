package client

import (
	"fmt"
)

// SavedSearchesService handles communication with the saved searches related
// methods of the Feedbin API
type SavedSearchesService struct {
	client *Client
}

// List returns all saved searches for the user
func (s *SavedSearchesService) List() ([]*SavedSearch, error) {
	req, err := s.client.newRequest("GET", "saved_searches.json", nil)
	if err != nil {
		return nil, err
	}

	var searches []*SavedSearch
	_, err = s.client.do(req, &searches)
	if err != nil {
		return nil, err
	}

	return searches, nil
}

// Get returns a single saved search
func (s *SavedSearchesService) Get(id int64) (*SavedSearch, error) {
	u := fmt.Sprintf("saved_searches/%d.json", id)

	req, err := s.client.newRequest("GET", u, nil)
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

// Create adds a new saved search
func (s *SavedSearchesService) Create(name, query string) (*SavedSearch, error) {
	type searchRequest struct {
		Name  string `json:"name"`
		Query string `json:"query"`
	}

	req, err := s.client.newRequest("POST", "saved_searches.json", &searchRequest{
		Name:  name,
		Query: query,
	})
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

// Update modifies a saved search
func (s *SavedSearchesService) Update(id int64, name, query string) (*SavedSearch, error) {
	type searchRequest struct {
		Name  string `json:"name,omitempty"`
		Query string `json:"query,omitempty"`
	}

	u := fmt.Sprintf("saved_searches/%d.json", id)

	req, err := s.client.newRequest("PATCH", u, &searchRequest{
		Name:  name,
		Query: query,
	})
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

// Delete removes a saved search
func (s *SavedSearchesService) Delete(id int64) error {
	u := fmt.Sprintf("saved_searches/%d.json", id)

	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	return err
}
