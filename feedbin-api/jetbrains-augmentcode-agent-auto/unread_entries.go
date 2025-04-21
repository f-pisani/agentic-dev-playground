package feedbin

import (
	"net/http"
)

// UnreadEntriesService handles communication with the unread entries related
// methods of the Feedbin API.
type UnreadEntriesService struct {
	client *Client
}

// List returns all unread entry IDs.
func (s *UnreadEntriesService) List() ([]int64, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/v2/unread_entries.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var entryIDs []int64
	resp, err := s.client.Do(req, &entryIDs)
	if err != nil {
		return nil, resp, err
	}

	return entryIDs, resp, nil
}

// UnreadEntriesRequest represents a request to mark entries as unread.
type UnreadEntriesRequest struct {
	UnreadEntries []int64 `json:"unread_entries"`
}

// Create marks entries as unread.
func (s *UnreadEntriesService) Create(entryIDs []int64) ([]int64, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/v2/unread_entries.json", &UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	})
	if err != nil {
		return nil, nil, err
	}

	var markedEntryIDs []int64
	resp, err := s.client.Do(req, &markedEntryIDs)
	if err != nil {
		return nil, resp, err
	}

	return markedEntryIDs, resp, nil
}

// Delete marks entries as read.
func (s *UnreadEntriesService) Delete(entryIDs []int64) ([]int64, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodDelete, "/v2/unread_entries.json", &UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	})
	if err != nil {
		return nil, nil, err
	}

	var markedEntryIDs []int64
	resp, err := s.client.Do(req, &markedEntryIDs)
	if err != nil {
		return nil, resp, err
	}

	return markedEntryIDs, resp, nil
}

// DeleteAlternative marks entries as read using the POST alternative for DELETE.
func (s *UnreadEntriesService) DeleteAlternative(entryIDs []int64) ([]int64, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/v2/unread_entries/delete.json", &UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	})
	if err != nil {
		return nil, nil, err
	}

	var markedEntryIDs []int64
	resp, err := s.client.Do(req, &markedEntryIDs)
	if err != nil {
		return nil, resp, err
	}

	return markedEntryIDs, resp, nil
}
