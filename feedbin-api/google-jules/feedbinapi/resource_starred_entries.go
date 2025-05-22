package feedbinapi

import (
	"context"
	"fmt"
	"net/http"
)

// StarredEntriesService handles operations related to starred entries.
type StarredEntriesService struct {
	client *Client
}

// ListStarredEntryIDs retrieves an array of entry IDs that are starred.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/starred-entries.md#get-starred-entries
func (s *StarredEntriesService) ListStarredEntryIDs(ctx context.Context) ([]int64, *Response, error) {
	path := "starred_entries.json"
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ListStarredEntryIDs request: %w", err)
	}

	var entryIDs []int64
	resp, err := s.client.do(req, &entryIDs)
	if err != nil {
		return nil, resp, err
	}
	return entryIDs, resp, nil
}

// StarEntries marks a list of entry IDs as starred.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/starred-entries.md#create-starred-entries
func (s *StarredEntriesService) StarEntries(ctx context.Context, entryIDs []int64) ([]int64, *Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, fmt.Errorf("entryIDs slice cannot be empty for StarEntries")
	}
	if len(entryIDs) > 1000 {
		// Note: API docs mention a limit of 1,000 entry_ids per request.
		// Consider returning an error or handling batching if this client should support more.
		// For now, assume the caller respects this limit or this is a simple client.
	}
	path := "starred_entries.json"
	requestBody := &StarredEntryRequest{StarredEntries: entryIDs}

	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create StarEntries request: %w", err)
	}

	var starredIDs []int64
	resp, err := s.client.do(req, &starredIDs)
	if err != nil {
		return nil, resp, err
	}
	return starredIDs, resp, nil
}

// UnstarEntries removes the starred status from a list of entry IDs.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/starred-entries.md#delete-starred-entries-unstar
func (s *StarredEntriesService) UnstarEntries(ctx context.Context, entryIDs []int64) ([]int64, *Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, fmt.Errorf("entryIDs slice cannot be empty for UnstarEntries")
	}
	path := "starred_entries.json"
	requestBody := &StarredEntryRequest{StarredEntries: entryIDs}

	req, err := s.client.newRequest(ctx, http.MethodDelete, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create UnstarEntries request: %w", err)
	}

	var unstarredIDs []int64
	resp, err := s.client.do(req, &unstarredIDs)
	if err != nil {
		return nil, resp, err
	}
	return unstarredIDs, resp, nil
}

// UnstarEntriesAlt provides the alternative POST method for unstarring entries.
// `POST /v2/starred_entries/delete.json`
func (s *StarredEntriesService) UnstarEntriesAlt(ctx context.Context, entryIDs []int64) ([]int64, *Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, fmt.Errorf("entryIDs slice cannot be empty for UnstarEntriesAlt")
	}
	path := "starred_entries/delete.json"
	requestBody := &StarredEntryRequest{StarredEntries: entryIDs}

	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create UnstarEntriesAlt request: %w", err)
	}
	var unstarredIDs []int64
	resp, err := s.client.do(req, &unstarredIDs)
	if err != nil {
		return nil, resp, err
	}
	return unstarredIDs, resp, nil
}
