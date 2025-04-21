package feedbin

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// SavedSearchesService handles communication with the saved searches related
// methods of the Feedbin API.
type SavedSearchesService struct {
	client *Client
}

// List returns all saved searches.
func (s *SavedSearchesService) List() ([]*SavedSearch, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/v2/saved_searches.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var savedSearches []*SavedSearch
	resp, err := s.client.Do(req, &savedSearches)
	if err != nil {
		return nil, resp, err
	}

	return savedSearches, resp, nil
}

// GetOptions specifies the optional parameters to the
// SavedSearchesService.Get method.
type GetOptions struct {
	IncludeEntries bool
	Page           int
}

// Get returns a single saved search.
func (s *SavedSearchesService) Get(id int64, opts *GetOptions) (interface{}, *http.Response, error) {
	u := fmt.Sprintf("/v2/saved_searches/%d.json", id)

	if opts != nil {
		params := url.Values{}

		if opts.IncludeEntries {
			params.Add("include_entries", "true")
		}

		if opts.Page > 0 {
			params.Add("page", strconv.Itoa(opts.Page))
		}

		if len(params) > 0 {
			u += "?" + params.Encode()
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	// The response can be either an array of entry IDs or an array of Entry objects
	// depending on the IncludeEntries parameter
	var result interface{}
	var resp *http.Response

	if opts != nil && opts.IncludeEntries {
		var entries []*Entry
		resp, err = s.client.Do(req, &entries)
		if err != nil {
			return nil, resp, err
		}
		result = entries
	} else {
		var entryIDs []int64
		resp, err = s.client.Do(req, &entryIDs)
		if err != nil {
			return nil, resp, err
		}
		result = entryIDs
	}

	return result, resp, nil
}

// CreateSavedSearchOptions specifies the parameters to the
// SavedSearchesService.Create method.
type CreateSavedSearchOptions struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// Create creates a new saved search.
func (s *SavedSearchesService) Create(opts *CreateSavedSearchOptions) (*SavedSearch, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/v2/saved_searches.json", opts)
	if err != nil {
		return nil, nil, err
	}

	savedSearch := new(SavedSearch)
	resp, err := s.client.Do(req, savedSearch)
	if err != nil {
		return nil, resp, err
	}

	return savedSearch, resp, nil
}

// Update updates a saved search.
func (s *SavedSearchesService) Update(id int64, opts *CreateSavedSearchOptions) (*SavedSearch, *http.Response, error) {
	u := fmt.Sprintf("/v2/saved_searches/%d.json", id)

	req, err := s.client.NewRequest(http.MethodPut, u, opts)
	if err != nil {
		return nil, nil, err
	}

	savedSearch := new(SavedSearch)
	resp, err := s.client.Do(req, savedSearch)
	if err != nil {
		return nil, resp, err
	}

	return savedSearch, resp, nil
}

// Delete deletes a saved search.
func (s *SavedSearchesService) Delete(id int64) (*http.Response, error) {
	u := fmt.Sprintf("/v2/saved_searches/%d.json", id)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
