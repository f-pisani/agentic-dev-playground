package feedbin

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// UnreadEntriesService handles communication with the unread entry related
// methods of the Feedbin API.
type UnreadEntriesService service

// unreadRequest is used for the body of MarkAsRead and MarkAsUnread requests.
type unreadRequest struct {
	UnreadEntries []int64 `json:"unread_entries"`
}

// List retrieves the IDs of all unread entries.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/unread-entries.md#get-unread-entries
func (s *UnreadEntriesService) List() ([]int64, *http.Response, error) {
	req, err := s.client.newRequest(http.MethodGet, "unread_entries.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var entryIDs []int64
	resp, err := s.client.do(req, &entryIDs)
	if err != nil {
		return nil, resp, err
	}

	return entryIDs, resp, nil
}

// MarkAsRead marks the specified entry IDs as read (i.e., removes them from the unread list).
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/unread-entries.md#delete-unread-entries
func (s *UnreadEntriesService) MarkAsRead(entryIDs []int64) (*http.Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil // Nothing to do
	}

	body := unreadRequest{UnreadEntries: entryIDs}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	req, err := s.client.newRequest(http.MethodDelete, "unread_entries.json", buf)
	if err != nil {
		return nil, err
	}

	// Expect 204 No Content
	resp, err := s.client.do(req, nil)
	return resp, err
}

// MarkAsUnread marks the specified entry IDs as unread.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/unread-entries.md#create-unread-entries
func (s *UnreadEntriesService) MarkAsUnread(entryIDs []int64) ([]int64, *http.Response, error) {
	if len(entryIDs) == 0 {
		return []int64{}, nil, nil // Nothing to do
	}

	body := unreadRequest{UnreadEntries: entryIDs}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest(http.MethodPost, "unread_entries.json", buf)
	if err != nil {
		return nil, nil, err
	}

	var createdIDs []int64 // API returns the IDs that were marked unread
	resp, err := s.client.do(req, &createdIDs)
	if err != nil {
		return nil, resp, err
	}

	return createdIDs, resp, nil
}
