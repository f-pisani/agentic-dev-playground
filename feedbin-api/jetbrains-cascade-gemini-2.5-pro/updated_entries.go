package feedbin

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// UpdatedEntriesService handles communication with the updated entry related
// methods of the Feedbin API.
type UpdatedEntriesService service

// updatedEntriesRequest is used for the body of the MarkAsRead request.
type updatedEntriesRequest struct {
	UpdatedEntries []int64 `json:"updated_entries"`
}

// List retrieves the IDs of entries that have been updated since the last request.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/updated-entries.md#get-updated-entries
func (s *UpdatedEntriesService) List() ([]int64, *http.Response, error) {
	req, err := s.client.newRequest(http.MethodGet, "updated_entries.json", nil)
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

// MarkAsRead marks the specified updated entry IDs as read (i.e., removes them from the updated list).
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/updated-entries.md#delete-updated-entries
func (s *UpdatedEntriesService) MarkAsRead(entryIDs []int64) (*http.Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil // Nothing to do
	}

	body := updatedEntriesRequest{UpdatedEntries: entryIDs}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	req, err := s.client.newRequest(http.MethodDelete, "updated_entries.json", buf)
	if err != nil {
		return nil, err
	}

	// Expect 204 No Content
	resp, err := s.client.do(req, nil)
	return resp, err
}
