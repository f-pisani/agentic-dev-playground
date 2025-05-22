package feedbinapi

import (
	"context"
	"fmt"
	"net/http"
)

// RecentlyReadEntriesService handles operations related to recently read entries.
type RecentlyReadEntriesService struct {
	client *Client
}

// ListRecentlyReadEntryIDs retrieves an array of entry IDs that were recently read.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/recently-read-entries.md#get-recently-read-entries
func (s *RecentlyReadEntriesService) ListRecentlyReadEntryIDs(ctx context.Context) ([]int64, *Response, error) {
	path := "recently_read_entries.json"
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ListRecentlyReadEntryIDs request: %w", err)
	}

	var entryIDs []int64
	resp, err := s.client.do(req, &entryIDs)
	if err != nil {
		return nil, resp, err
	}
	return entryIDs, resp, nil
}

// AddRecentlyReadEntries marks a list of entry IDs as recently read.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/recently-read-entries.md#create-recently-read-entries
func (s *RecentlyReadEntriesService) AddRecentlyReadEntries(ctx context.Context, entryIDs []int64) ([]int64, *Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, fmt.Errorf("entryIDs slice cannot be empty for AddRecentlyReadEntries")
	}
	path := "recently_read_entries.json"
	requestBody := &RecentlyReadEntryRequest{RecentlyReadEntries: entryIDs}

	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create AddRecentlyReadEntries request: %w", err)
	}

	var addedIDs []int64
	resp, err := s.client.do(req, &addedIDs)
	if err != nil {
		return nil, resp, err
	}
	return addedIDs, resp, nil
}
