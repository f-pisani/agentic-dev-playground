package feedbinapi

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// SavedSearchesService handles operations related to saved searches.
type SavedSearchesService struct {
	client *Client
}

// NewSavedSearchesService creates a new service for saved search operations.
func NewSavedSearchesService(client *Client) *SavedSearchesService {
	return &SavedSearchesService{client: client}
}

// SavedSearchListOptions specifies optional parameters for listing saved searches.
// The API doc for "GET /v2/saved_searches.json" does not specify query parameters.
type SavedSearchListOptions struct {
	ListOptions // Embeds Page, PerPage, Since (though not specified in docs, good for consistency)
}

// List retrieves all saved searches.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#get-saved-searches
func (s *SavedSearchesService) List(opts *SavedSearchListOptions) ([]SavedSearch, *http.Response, error) {
	path := "saved_searches.json"
	if opts != nil {
		v, err := query.Values(opts)
		if err != nil {
			return nil, nil, err
		}
		if params := v.Encode(); params != "" {
			path = fmt.Sprintf("%s?%s", path, params)
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var savedSearches []SavedSearch
	resp, err := s.client.Do(req, &savedSearches)
	if err != nil {
		return nil, resp, err
	}
	return savedSearches, resp, nil
}

// Get retrieves a single saved search by its ID.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#get-saved-search
func (s *SavedSearchesService) Get(id int64) (*SavedSearch, *http.Response, error) {
	if id == 0 {
		return nil, nil, fmt.Errorf("ID is required to get a saved search")
	}
	path := fmt.Sprintf("saved_searches/%d.json", id)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var savedSearch SavedSearch
	resp, err := s.client.Do(req, &savedSearch)
	if err != nil {
		return nil, resp, err
	}
	return &savedSearch, resp, nil
}

// CreateSavedSearchOptions specifies the parameters for creating a new saved search.
type CreateSavedSearchOptions struct {
	Name  string `json:"name"`
	Query string `json:"query"`
	// FeedIDs []int64 `json:"feed_ids,omitempty"` // Optional: according to model, might not be part of create
}

// Create creates a new saved search.
// Request body: {"name": "Search Name", "query": "search query"}
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#create-saved-search
func (s *SavedSearchesService) Create(opts *CreateSavedSearchOptions) (*SavedSearch, *http.Response, error) {
	if opts == nil || opts.Name == "" || opts.Query == "" {
		return nil, nil, fmt.Errorf("Name and Query are required to create a saved search")
	}
	path := "saved_searches.json"

	req, err := s.client.NewRequest(http.MethodPost, path, opts)
	if err != nil {
		return nil, nil, err
	}

	var createdSearch SavedSearch
	resp, err := s.client.Do(req, &createdSearch)
	if err != nil {
		return nil, resp, err
	}
	return &createdSearch, resp, nil
}

// UpdateSavedSearchOptions specifies the parameters for updating a saved search.
type UpdateSavedSearchOptions struct {
	Name  string `json:"name,omitempty"`  // If not provided, will not be updated
	Query string `json:"query,omitempty"` // If not provided, will not be updated
	// FeedIDs []int64 `json:"feed_ids,omitempty"` // Optional
}

// Update modifies an existing saved search.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#update-saved-search
func (s *SavedSearchesService) Update(id int64, opts *UpdateSavedSearchOptions) (*SavedSearch, *http.Response, error) {
	if id == 0 {
		return nil, nil, fmt.Errorf("ID is required to update a saved search")
	}
	if opts == nil || (opts.Name == "" && opts.Query == "") { // At least one field must be present for update
		return nil, nil, fmt.Errorf("Either Name or Query must be provided to update a saved search")
	}
	path := fmt.Sprintf("saved_searches/%d.json", id)

	req, err := s.client.NewRequest(http.MethodPatch, path, opts) // PATCH for partial updates
	if err != nil {
		return nil, nil, err
	}

	var updatedSearch SavedSearch
	resp, err := s.client.Do(req, &updatedSearch)
	if err != nil {
		return nil, resp, err
	}
	return &updatedSearch, resp, nil
}

// Delete removes a saved search by its ID.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#delete-saved-search
func (s *SavedSearchesService) Delete(id int64) (*http.Response, error) {
	if id == 0 {
		return nil, fmt.Errorf("ID is required to delete a saved search")
	}
	path := fmt.Sprintf("saved_searches/%d.json", id)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil) // Expect 204 No Content
	if err != nil {
		return resp, err
	}
	return resp, nil
}
