package feedbin

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// RecentlyReadService handles communication with the recently read entries
// related methods of the Feedbin API.
type RecentlyReadService service

// List retrieves the IDs of recently read entries.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/recently-read-entries.md#get-recently-read-entries
func (s *RecentlyReadService) List() ([]int64, *http.Response, error) {
	req, err := s.client.newRequest(http.MethodGet, "recently_read_entries.json", nil)
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

// recentlyReadRequest is used for the body of MarkAsRead requests.
type recentlyReadRequest struct {
	RecentlyReadEntries []int64 `json:"recently_read_entries"`
}

// MarkAsRead marks the specified entry IDs as recently read.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/recently-read-entries.md#create-recently-read-entries
func (s *RecentlyReadService) MarkAsRead(entryIDs []int64) ([]int64, *http.Response, error) {
	if len(entryIDs) == 0 {
		return []int64{}, nil, nil // Nothing to do
	}

	body := recentlyReadRequest{RecentlyReadEntries: entryIDs}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest(http.MethodPost, "recently_read_entries.json", buf)
	if err != nil {
		return nil, nil, err
	}

	var createdIDs []int64
	resp, err := s.client.do(req, &createdIDs)
	if err != nil {
		return nil, resp, err
	}

	return createdIDs, resp, nil
}
