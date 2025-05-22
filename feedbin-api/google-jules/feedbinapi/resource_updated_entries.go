package feedbinapi

import (
	"context"
	"fmt"
	"net/http"
)

// UpdatedEntriesService handles operations related to entries that have been updated since publication.
type UpdatedEntriesService struct {
	client *Client
}

// ListUpdatedEntryIDs retrieves an array of entry IDs that have been updated.
// Supports 'since' parameter.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/updated-entries.md#get-updated-entries
func (s *UpdatedEntriesService) ListUpdatedEntryIDs(ctx context.Context, params *ListUpdatedEntryIDsParams) ([]int64, *Response, error) {
	path := "updated_entries.json"

	queryParams, err := structToURLValues(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert params to query values for ListUpdatedEntryIDs: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ListUpdatedEntryIDs request: %w", err)
	}

	var entryIDs []int64
	resp, err := s.client.do(req, &entryIDs)
	if err != nil {
		return nil, resp, err
	}
	return entryIDs, resp, nil
}

// MarkUpdatesAsRead marks specified updated entry IDs as "read" (acknowledged the update).
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/updated-entries.md#delete-updated-entries-mark-as-read
func (s *UpdatedEntriesService) MarkUpdatesAsRead(ctx context.Context, entryIDs []int64) ([]int64, *Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, fmt.Errorf("entryIDs slice cannot be empty for MarkUpdatesAsRead")
	}
	path := "updated_entries.json"
	requestBody := &UpdatedEntryRequest{UpdatedEntries: entryIDs}

	req, err := s.client.newRequest(ctx, http.MethodDelete, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create MarkUpdatesAsRead request: %w", err)
	}

	var processedIDs []int64
	resp, err := s.client.do(req, &processedIDs)
	if err != nil {
		return nil, resp, err
	}
	return processedIDs, resp, nil
}

// MarkUpdatesAsReadAlt provides the alternative POST method for marking updated entries as read.
// `POST /v2/updated_entries/delete.json`
func (s *UpdatedEntriesService) MarkUpdatesAsReadAlt(ctx context.Context, entryIDs []int64) ([]int64, *Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, fmt.Errorf("entryIDs slice cannot be empty for MarkUpdatesAsReadAlt")
	}
	path := "updated_entries/delete.json"
	requestBody := &UpdatedEntryRequest{UpdatedEntries: entryIDs}

	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create MarkUpdatesAsReadAlt request: %w", err)
	}
	var processedIDs []int64
	resp, err := s.client.do(req, &processedIDs)
	if err != nil {
		return nil, resp, err
	}
	return processedIDs, resp, nil
}
