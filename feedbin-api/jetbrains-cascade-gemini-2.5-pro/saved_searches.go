package feedbin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SavedSearchesService handles communication with the saved search related
// methods of the Feedbin API.
type SavedSearchesService service

// List retrieves all saved searches for the user.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#get-saved-searches
func (s *SavedSearchesService) List() ([]*SavedSearch, *http.Response, error) {
	req, err := s.client.newRequest(http.MethodGet, "saved_searches.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var searches []*SavedSearch
	resp, err := s.client.do(req, &searches)
	if err != nil {
		return nil, resp, err
	}

	return searches, resp, nil
}

// GetSavedSearchEntriesOptions specifies the optional parameters for GetEntries.
type GetSavedSearchEntriesOptions struct {
	IncludeEntries *bool `url:"include_entries,omitempty"`
	Page           *int  `url:"page,omitempty"`
}

// GetEntriesResponse defines the possible response structures for GetEntries.
// It can be either a slice of entry IDs or a slice of full Entry objects.
type GetEntriesResponse struct {
	EntryIDs []int64
	Entries  []*Entry
}

// UnmarshalJSON handles the custom unmarshalling logic.
func (r *GetEntriesResponse) UnmarshalJSON(data []byte) error {
	// Try unmarshalling into []int64 first
	if err := json.Unmarshal(data, &r.EntryIDs); err == nil {
		return nil // Success, it's an array of IDs
	}

	// If that fails, try unmarshalling into []*Entry
	if err := json.Unmarshal(data, &r.Entries); err == nil {
		return nil // Success, it's an array of Entry objects
	}

	// If both fail, return an error
	return fmt.Errorf("failed to unmarshal saved search entries response: %s", string(data))
}

// GetEntries retrieves the entries matching a specific saved search.
// By default, it returns an array of entry IDs.
// Use opt.IncludeEntries = true to get full Entry objects.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#get-saved-search
func (s *SavedSearchesService) GetEntries(searchID int64, opt *GetSavedSearchEntriesOptions) (*GetEntriesResponse, *http.Response, error) {
	path := fmt.Sprintf("saved_searches/%d.json", searchID)
	u, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result GetEntriesResponse
	resp, err := s.client.do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// savedSearchRequest is used for Create and Update requests.
type savedSearchRequest struct {
	Name  *string `json:"name,omitempty"`
	Query *string `json:"query,omitempty"`
}

// Create creates a new saved search.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#create-saved-search
func (s *SavedSearchesService) Create(name, query string) (*SavedSearch, *http.Response, error) {
	body := savedSearchRequest{
		Name:  &name,
		Query: &query,
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest(http.MethodPost, "saved_searches.json", buf)
	if err != nil {
		return nil, nil, err
	}

	// Expect 201 Created, response body is the created search, but we don't decode it here.
	// The Location header contains the URL to the new resource.
	var createdSearch SavedSearch // To potentially get ID if returned in body (though docs say Location header)
	resp, err := s.client.do(req, &createdSearch)
	if err != nil {
		return nil, resp, err
	}

	// It's safer to fetch the newly created search using the Location header if present.
	// However, the API might return the created object directly in the 201 response body.
	return &createdSearch, resp, nil
}

// Update modifies an existing saved search.
// Only the Name can be updated according to the docs for PATCH.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#update-saved-search
func (s *SavedSearchesService) Update(searchID int64, newName string) (*SavedSearch, *http.Response, error) {
	path := fmt.Sprintf("saved_searches/%d.json", searchID)
	body := savedSearchRequest{Name: &newName}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest(http.MethodPatch, path, buf)
	if err != nil {
		return nil, nil, err
	}

	var updatedSearch SavedSearch
	resp, err := s.client.do(req, &updatedSearch)
	if err != nil {
		return nil, resp, err
	}

	return &updatedSearch, resp, nil
}

// UpdateWithPost modifies an existing saved search using the POST alternative endpoint.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#patch-alternative
func (s *SavedSearchesService) UpdateWithPost(searchID int64, newName string) (*SavedSearch, *http.Response, error) {
	path := fmt.Sprintf("saved_searches/%d/update.json", searchID)
	body := savedSearchRequest{Name: &newName}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest(http.MethodPost, path, buf)
	if err != nil {
		return nil, nil, err
	}

	var updatedSearch SavedSearch
	resp, err := s.client.do(req, &updatedSearch)
	if err != nil {
		return nil, resp, err
	}

	return &updatedSearch, resp, nil
}

// Delete removes a saved search.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/saved-searches.md#delete-saved-search
func (s *SavedSearchesService) Delete(searchID int64) (*http.Response, error) {
	path := fmt.Sprintf("saved_searches/%d.json", searchID)
	req, err := s.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	// Expect 204 No Content
	resp, err := s.client.do(req, nil)
	return resp, err
}
