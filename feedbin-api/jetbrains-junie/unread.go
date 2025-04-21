// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"net/http"
)

// UnreadService provides access to the unread entries endpoints of the Feedbin API.
type UnreadService struct {
	client *Client
}

// UnreadEntriesRequest represents a request to mark entries as read or unread.
type UnreadEntriesRequest struct {
	UnreadEntries []int `json:"unread_entries"`
}

// GetUnreadEntries retrieves all unread entry IDs.
func (s *UnreadService) GetUnreadEntries() ([]int, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "/unread_entries.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var entryIDs []int
	resp, err := s.client.Do(req, &entryIDs)
	if err != nil {
		return nil, resp, err
	}

	return entryIDs, resp, nil
}

// MarkAsUnread marks the specified entries as unread.
func (s *UnreadService) MarkAsUnread(entryIDs []int) ([]int, *http.Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, nil
	}

	if len(entryIDs) > 1000 {
		return nil, nil, &APIError{
			Message: "maximum of 1,000 entry IDs allowed",
		}
	}

	req, err := s.client.NewRequest("POST", "/unread_entries.json", &UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	})
	if err != nil {
		return nil, nil, err
	}

	var markedEntryIDs []int
	resp, err := s.client.Do(req, &markedEntryIDs)
	if err != nil {
		return nil, resp, err
	}

	return markedEntryIDs, resp, nil
}

// MarkAsRead marks the specified entries as read.
func (s *UnreadService) MarkAsRead(entryIDs []int) ([]int, *http.Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, nil
	}

	if len(entryIDs) > 1000 {
		return nil, nil, &APIError{
			Message: "maximum of 1,000 entry IDs allowed",
		}
	}

	req, err := s.client.NewRequest("DELETE", "/unread_entries.json", &UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	})
	if err != nil {
		return nil, nil, err
	}

	var markedEntryIDs []int
	resp, err := s.client.Do(req, &markedEntryIDs)
	if err != nil {
		return nil, resp, err
	}

	return markedEntryIDs, resp, nil
}

// MarkAsReadAlternative marks the specified entries as read using the POST alternative.
// Some clients like Android don't easily allow a body with a DELETE request.
func (s *UnreadService) MarkAsReadAlternative(entryIDs []int) ([]int, *http.Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, nil
	}

	if len(entryIDs) > 1000 {
		return nil, nil, &APIError{
			Message: "maximum of 1,000 entry IDs allowed",
		}
	}

	req, err := s.client.NewRequest("POST", "/unread_entries/delete.json", &UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	})
	if err != nil {
		return nil, nil, err
	}

	var markedEntryIDs []int
	resp, err := s.client.Do(req, &markedEntryIDs)
	if err != nil {
		return nil, resp, err
	}

	return markedEntryIDs, resp, nil
}
