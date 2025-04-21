package feedbin

import (
	"fmt"
	"net/http"
	neturl "net/url"
	"strconv"
)

// SavedSearchesService handles communication with the saved searches related
// endpoints of the Feedbin API
type SavedSearchesService struct {
	client *Client
}

// SavedSearchCreateOptions specifies the parameters to the
// SavedSearchesService.Create method
type SavedSearchCreateOptions struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// SavedSearchUpdateOptions specifies the parameters to the
// SavedSearchesService.Update method
type SavedSearchUpdateOptions struct {
	Name  string `json:"name,omitempty"`
	Query string `json:"query,omitempty"`
}

// SavedSearchGetOptions specifies the optional parameters to the
// SavedSearchesService.Get method
type SavedSearchGetOptions struct {
	IncludeEntries bool
	Page           int
}

// List returns all saved searches for the authenticated user
func (s *SavedSearchesService) List() ([]*SavedSearch, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "saved_searches.json", nil)
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

// Get retrieves the entries matching a saved search
func (s *SavedSearchesService) Get(id int, opts *SavedSearchGetOptions) (interface{}, *PaginationInfo, *http.Response, error) {
	url := fmt.Sprintf("saved_searches/%d.json", id)
	params := neturl.Values{}

	if opts != nil {
		if opts.IncludeEntries {
			params.Add("include_entries", "true")
		}
		if opts.Page > 0 {
			params.Add("page", strconv.Itoa(opts.Page))
		}
	}

	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	if opts != nil && opts.IncludeEntries {
		var entries []*Entry
		resp, err := s.client.Do(req, &entries)
		if err != nil {
			return nil, nil, resp, err
		}
		pagination := s.client.GetPagination(resp)
		return entries, pagination, resp, nil
	} else {
		var entryIDs []int
		resp, err := s.client.Do(req, &entryIDs)
		if err != nil {
			return nil, nil, resp, err
		}
		pagination := s.client.GetPagination(resp)
		return entryIDs, pagination, resp, nil
	}
}

// Create creates a new saved search
func (s *SavedSearchesService) Create(opts *SavedSearchCreateOptions) (*SavedSearch, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "saved_searches.json", opts)
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

// Delete deletes a saved search
func (s *SavedSearchesService) Delete(id int) (*http.Response, error) {
	url := fmt.Sprintf("saved_searches/%d.json", id)
	req, err := s.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Update updates a saved search
func (s *SavedSearchesService) Update(id int, opts *SavedSearchUpdateOptions) (*SavedSearch, *http.Response, error) {
	url := fmt.Sprintf("saved_searches/%d.json", id)
	req, err := s.client.NewRequest(http.MethodPatch, url, opts)
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

// UpdateWithPost updates a saved search using POST instead of PATCH
// Some proxies may block PATCH requests
func (s *SavedSearchesService) UpdateWithPost(id int, opts *SavedSearchUpdateOptions) (*SavedSearch, *http.Response, error) {
	url := fmt.Sprintf("saved_searches/%d/update.json", id)
	req, err := s.client.NewRequest(http.MethodPost, url, opts)
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
