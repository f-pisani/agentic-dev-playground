package feedbin

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// SavedSearchesService handles operations related to saved searches
type SavedSearchesService struct {
	client *Client
}

// List returns all saved searches
func (s *SavedSearchesService) List() ([]SavedSearch, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/saved_searches.json", nil)
	if err != nil {
		return nil, err
	}

	var savedSearches []SavedSearch
	_, err = s.client.Do(req, &savedSearches)
	if err != nil {
		return nil, err
	}

	return savedSearches, nil
}

// Get returns the entry IDs matching a saved search
// If includeEntries is true, returns the entries instead of entry IDs
// If page is provided, returns the specified page of results
func (s *SavedSearchesService) Get(id int, includeEntries bool, page *int) (interface{}, error) {
	path := fmt.Sprintf("/saved_searches/%d.json", id)

	// Add query parameters if needed
	values := url.Values{}
	if includeEntries {
		values.Add("include_entries", "true")
	}
	if page != nil {
		values.Add("page", strconv.Itoa(*page))
	}

	if len(values) > 0 {
		path = path + "?" + values.Encode()
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// If includeEntries is true, we expect entries, otherwise we expect entry IDs
	if includeEntries {
		var entries []Entry
		_, err = s.client.Do(req, &entries)
		if err != nil {
			return nil, err
		}
		return entries, nil
	} else {
		var entryIDs []int
		_, err = s.client.Do(req, &entryIDs)
		if err != nil {
			return nil, err
		}
		return entryIDs, nil
	}
}

// Create creates a new saved search
func (s *SavedSearchesService) Create(name, query string) (*SavedSearch, error) {
	params := &SavedSearchParams{
		Name:  name,
		Query: query,
	}

	req, err := s.client.NewRequest(http.MethodPost, "/saved_searches.json", params)
	if err != nil {
		return nil, err
	}

	savedSearch := new(SavedSearch)
	_, err = s.client.Do(req, savedSearch)
	if err != nil {
		return nil, err
	}

	return savedSearch, nil
}

// Update updates a saved search
func (s *SavedSearchesService) Update(id int, name string) (*SavedSearch, error) {
	path := fmt.Sprintf("/saved_searches/%d.json", id)

	params := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}

	req, err := s.client.NewRequest(http.MethodPatch, path, params)
	if err != nil {
		return nil, err
	}

	savedSearch := new(SavedSearch)
	_, err = s.client.Do(req, savedSearch)
	if err != nil {
		return nil, err
	}

	return savedSearch, nil
}

// UpdateWithPost is an alternative to Update that uses POST instead of PATCH
// This is useful for clients that don't support PATCH requests
func (s *SavedSearchesService) UpdateWithPost(id int, name string) (*SavedSearch, error) {
	path := fmt.Sprintf("/saved_searches/%d/update.json", id)

	params := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}

	req, err := s.client.NewRequest(http.MethodPost, path, params)
	if err != nil {
		return nil, err
	}

	savedSearch := new(SavedSearch)
	_, err = s.client.Do(req, savedSearch)
	if err != nil {
		return nil, err
	}

	return savedSearch, nil
}

// Delete deletes a saved search
func (s *SavedSearchesService) Delete(id int) error {
	path := fmt.Sprintf("/saved_searches/%d.json", id)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}
