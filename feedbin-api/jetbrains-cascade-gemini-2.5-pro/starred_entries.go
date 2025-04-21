package feedbin

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// StarredEntriesService handles communication with the starred entry related
// methods of the Feedbin API.
type StarredEntriesService service

// starredRequest is used for the body of MarkAsStarred and MarkAsUnstarred requests.
type starredRequest struct {
	StarredEntries []int64 `json:"starred_entries"`
}

// List retrieves the IDs of all starred entries.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/starred-entries.md#get-starred-entries
func (s *StarredEntriesService) List() ([]int64, *http.Response, error) {
	req, err := s.client.newRequest(http.MethodGet, "starred_entries.json", nil)
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

// MarkAsStarred marks the specified entry IDs as starred.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/starred-entries.md#create-starred-entries
func (s *StarredEntriesService) MarkAsStarred(entryIDs []int64) ([]int64, *http.Response, error) {
	if len(entryIDs) == 0 {
		return []int64{}, nil, nil // Nothing to do
	}

	body := starredRequest{StarredEntries: entryIDs}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest(http.MethodPost, "starred_entries.json", buf)
	if err != nil {
		return nil, nil, err
	}

	var createdIDs []int64 // API returns the IDs that were starred
	resp, err := s.client.do(req, &createdIDs)
	if err != nil {
		return nil, resp, err
	}

	return createdIDs, resp, nil
}

// MarkAsUnstarred marks the specified entry IDs as unstarred.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/starred-entries.md#delete-starred-entries
func (s *StarredEntriesService) MarkAsUnstarred(entryIDs []int64) (*http.Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil // Nothing to do
	}

	body := starredRequest{StarredEntries: entryIDs}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	req, err := s.client.newRequest(http.MethodDelete, "starred_entries.json", buf)
	if err != nil {
		return nil, err
	}

	// Expect 204 No Content
	resp, err := s.client.do(req, nil)
	return resp, err
}
