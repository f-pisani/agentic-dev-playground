package feedbinapi

import (
	"context"
	"fmt"
	"net/http"
)

// UnreadEntriesService handles operations related to unread entries.
type UnreadEntriesService struct {
	client *Client
}

// ListUnreadEntryIDs retrieves an array of entry IDs that are unread.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/unread-entries.md#get-unread-entries
func (s *UnreadEntriesService) ListUnreadEntryIDs(ctx context.Context) ([]int64, *Response, error) {
	path := "unread_entries.json"
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ListUnreadEntryIDs request: %w", err)
	}

	var entryIDs []int64
	resp, err := s.client.do(req, &entryIDs)
	if err != nil {
		return nil, resp, err
	}
	return entryIDs, resp, nil
}

// MarkAsUnread marks a list of entry IDs as unread.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/unread-entries.md#create-unread-entries-mark-as-unread
func (s *UnreadEntriesService) MarkAsUnread(ctx context.Context, entryIDs []int64) ([]int64, *Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, fmt.Errorf("entryIDs slice cannot be empty for MarkAsUnread")
	}
	path := "unread_entries.json"
	requestBody := &UnreadEntryRequest{UnreadEntries: entryIDs}

	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create MarkAsUnread request: %w", err)
	}

	var markedIDs []int64
	resp, err := s.client.do(req, &markedIDs)
	if err != nil {
		return nil, resp, err
	}
	return markedIDs, resp, nil
}

// MarkAsRead marks a list of entry IDs as read.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/unread-entries.md#delete-unread-entries-mark-as-read
func (s *UnreadEntriesService) MarkAsRead(ctx context.Context, entryIDs []int64) ([]int64, *Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, fmt.Errorf("entryIDs slice cannot be empty for MarkAsRead")
	}
	path := "unread_entries.json"
	requestBody := &UnreadEntryRequest{UnreadEntries: entryIDs}

	req, err := s.client.newRequest(ctx, http.MethodDelete, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create MarkAsRead request: %w", err)
	}

	var markedIDs []int64
	resp, err := s.client.do(req, &markedIDs)
	if err != nil {
		return nil, resp, err
	}
	return markedIDs, resp, nil
}

// MarkAsReadAlt provides the alternative POST method for marking entries as read.
// `POST /v2/unread_entries/delete.json`
func (s *UnreadEntriesService) MarkAsReadAlt(ctx context.Context, entryIDs []int64) ([]int64, *Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, fmt.Errorf("entryIDs slice cannot be empty for MarkAsReadAlt")
	}
	path := "unread_entries/delete.json"
	requestBody := &UnreadEntryRequest{UnreadEntries: entryIDs}

	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create MarkAsReadAlt request: %w", err)
	}
	var markedIDs []int64
	resp, err := s.client.do(req, &markedIDs)
	if err != nil {
		return nil, resp, err
	}
	return markedIDs, resp, nil
}
