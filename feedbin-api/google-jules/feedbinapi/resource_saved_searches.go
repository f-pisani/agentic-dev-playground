package feedbinapi

import (
	"context"
	"fmt"
	"net/http"
	// "encoding/json" // Potentially needed for GetSavedSearch if it returns mixed types or requires manual parsing for entry IDs vs full entries
)

// SavedSearchesService handles operations related to saved searches.
type SavedSearchesService struct {
	client *Client
}

// ListSavedSearches retrieves all saved searches for the authenticated user.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/saved-searches.md#get-saved-searches
func (s *SavedSearchesService) ListSavedSearches(ctx context.Context) ([]SavedSearch, *Response, error) {
	path := "saved_searches.json"
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ListSavedSearches request: %w", err)
	}

	var searches []SavedSearch
	resp, err := s.client.do(req, &searches)
	if err != nil {
		return nil, resp, err
	}
	return searches, resp, nil
}

// GetSavedSearchEntryIDs retrieves entry IDs for a specific saved search.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/saved-searches.md#get-saved-search
// (By default, an array of entry ids is returned)
func (s *SavedSearchesService) GetSavedSearchEntryIDs(ctx context.Context, savedSearchID int64) ([]int64, *Response, error) {
	if savedSearchID <= 0 {
		return nil, nil, fmt.Errorf("savedSearchID must be a positive integer")
	}
	path := fmt.Sprintf("saved_searches/%d.json", savedSearchID)
	// No params, specifically `include_entries` is false by default

	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GetSavedSearchEntryIDs request: %w", err)
	}

	var entryIDs []int64
	resp, err := s.client.do(req, &entryIDs) // Expects the response to be a simple JSON array of numbers
	if err != nil {
		return nil, resp, err
	}
	return entryIDs, resp, nil
}

// GetSavedSearchWithEntries retrieves full entry objects for a specific saved search.
// This is paginated.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/saved-searches.md#get-saved-search
// (GET /v2/saved_searches/1.json?include_entries=true will return entry objects)
func (s *SavedSearchesService) GetSavedSearchWithEntries(ctx context.Context, savedSearchID int64, params *GetSavedSearchParams) ([]Entry, *Response, error) {
	if savedSearchID <= 0 {
		return nil, nil, fmt.Errorf("savedSearchID must be a positive integer")
	}
	path := fmt.Sprintf("saved_searches/%d.json", savedSearchID)

	var effectiveParams GetSavedSearchParams
	if params != nil {
		effectiveParams = *params
	}
	effectiveParams.IncludeEntries = boolPtr(true) // Ensure IncludeEntries is true

	queryParams, err := structToURLValues(&effectiveParams)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert params for GetSavedSearchWithEntries: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GetSavedSearchWithEntries request: %w", err)
	}

	var entries []Entry // Expects the response to be a JSON array of Entry objects
	resp, err := s.client.do(req, &entries)
	if err != nil {
		return nil, resp, err
	}
	return entries, resp, nil
}

// CreateSavedSearch creates a new saved search.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/saved-searches.md#create-saved-search
func (s *SavedSearchesService) CreateSavedSearch(ctx context.Context, name string, query string) (*SavedSearch, *Response, error) {
	if name == "" || query == "" {
		return nil, nil, fmt.Errorf("name and query cannot be empty for CreateSavedSearch")
	}
	path := "saved_searches.json"
	requestBody := &CreateSavedSearchRequest{Name: name, Query: query}

	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CreateSavedSearch request: %w", err)
	}

	var createdSearch SavedSearch
	resp, err := s.client.do(req, &createdSearch)
	if err != nil {
		return nil, resp, err
	}
	return &createdSearch, resp, nil
}

// DeleteSavedSearch deletes a saved search by its ID.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/saved-searches.md#delete-saved-search
func (s *SavedSearchesService) DeleteSavedSearch(ctx context.Context, savedSearchID int64) (*Response, error) {
	if savedSearchID <= 0 {
		return nil, fmt.Errorf("savedSearchID must be a positive integer")
	}
	path := fmt.Sprintf("saved_searches/%d.json", savedSearchID)

	req, err := s.client.newRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DeleteSavedSearch request: %w", err)
	}

	resp, err := s.client.do(req, nil) // Expects 204 No Content
	return resp, err
}

// UpdateSavedSearch updates a saved search, typically its name.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/saved-searches.md#update-saved-search
func (s *SavedSearchesService) UpdateSavedSearch(ctx context.Context, savedSearchID int64, name string) (*SavedSearch, *Response, error) {
	if savedSearchID <= 0 {
		return nil, nil, fmt.Errorf("savedSearchID must be a positive integer")
	}
	if name == "" {
		return nil, nil, fmt.Errorf("name cannot be empty for UpdateSavedSearch")
	}
	path := fmt.Sprintf("saved_searches/%d.json", savedSearchID)
	requestBody := &UpdateSavedSearchRequest{Name: &name}

	req, err := s.client.newRequest(ctx, http.MethodPatch, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create UpdateSavedSearch (PATCH) request: %w", err)
	}

	var updatedSearch SavedSearch
	resp, err := s.client.do(req, &updatedSearch)
	if err != nil {
		return nil, resp, err
	}
	return &updatedSearch, resp, nil
}

// UpdateSavedSearchAlt handles updating a saved search using the POST alternative.
// `POST /v2/saved_searches/:id/update.json`
func (s *SavedSearchesService) UpdateSavedSearchAlt(ctx context.Context, savedSearchID int64, name string) (*SavedSearch, *Response, error) {
	if savedSearchID <= 0 {
		return nil, nil, fmt.Errorf("savedSearchID must be a positive integer")
	}
	if name == "" {
		return nil, nil, fmt.Errorf("name cannot be empty")
	}
	path := fmt.Sprintf("saved_searches/%d/update.json", savedSearchID)
	requestBody := &UpdateSavedSearchRequest{Name: &name}

	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create UpdateSavedSearch (POST alt) request: %w", err)
	}
	var updatedSearch SavedSearch
	resp, err := s.client.do(req, &updatedSearch)
	if err != nil {
		return nil, resp, err
	}
	return &updatedSearch, resp, nil
}

// Helper for pointers to bool, used in GetSavedSearchWithEntries
func boolPtr(b bool) *bool { return &b }

// The generic GetSavedSearch is removed in favor of GetSavedSearchEntryIDs and GetSavedSearchWithEntries
// to handle the API's response type change based on `include_entries` parameter.
/*
func (s *SavedSearchesService) GetSavedSearch(ctx context.Context, savedSearchID int64, params *GetSavedSearchParams) (*SavedSearchDetail, *Response, error) {
	if savedSearchID <= 0 {
		return nil, nil, fmt.Errorf("savedSearchID must be a positive integer")
	}
	path := fmt.Sprintf("saved_searches/%d.json", savedSearchID)

	queryParams, err := structToURLValues(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert params to query values for GetSavedSearch: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GetSavedSearch request: %w", err)
	}

	var detail SavedSearchDetail // Requires SavedSearchDetail to have custom unmarshalling
	resp, err := s.client.do(req, &detail)
	if err != nil {
		return nil, resp, err
	}
	return &detail, resp, nil
}
*/
