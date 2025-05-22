package feedbinapi

import (
	"context"
	"fmt"
	"net/http"
)

// EntriesService handles operations related to entries.
type EntriesService struct {
	client *Client
}

// ListEntries retrieves a list of entries for the authenticated user.
// It supports pagination and various filtering options via ListEntriesParams.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/entries.md#get-v2entriesjson
func (s *EntriesService) ListEntries(ctx context.Context, params *ListEntriesParams) ([]Entry, *Response, error) {
	path := "entries.json"
	
	queryParams, err := structToURLValues(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert params to query values for ListEntries: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ListEntries request: %w", err)
	}

	var entries []Entry
	resp, err := s.client.do(req, &entries)
	if err != nil {
		return nil, resp, err
	}

	return entries, resp, nil
}

// GetEntry retrieves a single entry by its ID.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/entries.md#get-v2entries3648json
func (s *EntriesService) GetEntry(ctx context.Context, entryID int64, params *GetEntryParams) (*Entry, *Response, error) {
	if entryID <= 0 {
		return nil, nil, fmt.Errorf("entryID must be a positive integer")
	}
	path := fmt.Sprintf("entries/%d.json", entryID)

	queryParams, err := structToURLValues(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert params to query values for GetEntry: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GetEntry request: %w", err)
	}

	var entry Entry
	resp, err := s.client.do(req, &entry)
	if err != nil {
		return nil, resp, err
	}

	return &entry, resp, nil
}

// ListFeedEntries retrieves a list of entries for a specific feed.
// Supports pagination and filtering options via ListEntriesParams (excluding 'ids').
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/entries.md#get-v2feeds203entriesjson
func (s *EntriesService) ListFeedEntries(ctx context.Context, feedID int64, params *ListEntriesParams) ([]Entry, *Response, error) {
	if feedID <= 0 {
		return nil, nil, fmt.Errorf("feedID must be a positive integer")
	}
	path := fmt.Sprintf("feeds/%d/entries.json", feedID)

	// Ensure 'ids' param is not used for this endpoint, as per docs
	var effectiveParams *ListEntriesParams
	if params != nil {
		pCopy := *params
		pCopy.IDs = nil // IDs parameter is not supported for this endpoint
		effectiveParams = &pCopy
	}

	queryParams, err := structToURLValues(effectiveParams)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert params to query values for ListFeedEntries: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ListFeedEntries request: %w", err)
	}

	var entries []Entry
	resp, err := s.client.do(req, &entries)
	if err != nil {
		return nil, resp, err
	}

	return entries, resp, nil
}
