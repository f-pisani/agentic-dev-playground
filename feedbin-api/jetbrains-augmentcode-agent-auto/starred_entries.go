package feedbin

import (
	"net/http"
)

// StarredEntriesService handles communication with the starred entries related
// methods of the Feedbin API.
type StarredEntriesService struct {
	client *Client
}

// List returns all starred entry IDs.
func (s *StarredEntriesService) List() ([]int64, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/v2/starred_entries.json", nil)
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

// StarredEntriesRequest represents a request to star entries.
type StarredEntriesRequest struct {
	StarredEntries []int64 `json:"starred_entries"`
}

// Create stars entries.
func (s *StarredEntriesService) Create(entryIDs []int64) ([]int64, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/v2/starred_entries.json", &StarredEntriesRequest{
		StarredEntries: entryIDs,
	})
	if err != nil {
		return nil, nil, err
	}

	var starredEntryIDs []int64
	resp, err := s.client.Do(req, &starredEntryIDs)
	if err != nil {
		return nil, resp, err
	}

	return starredEntryIDs, resp, nil
}

// Delete unstars entries.
func (s *StarredEntriesService) Delete(entryIDs []int64) ([]int64, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodDelete, "/v2/starred_entries.json", &StarredEntriesRequest{
		StarredEntries: entryIDs,
	})
	if err != nil {
		return nil, nil, err
	}

	var unstarredEntryIDs []int64
	resp, err := s.client.Do(req, &unstarredEntryIDs)
	if err != nil {
		return nil, resp, err
	}

	return unstarredEntryIDs, resp, nil
}

// DeleteAlternative unstars entries using the POST alternative for DELETE.
func (s *StarredEntriesService) DeleteAlternative(entryIDs []int64) ([]int64, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/v2/starred_entries/delete.json", &StarredEntriesRequest{
		StarredEntries: entryIDs,
	})
	if err != nil {
		return nil, nil, err
	}

	var unstarredEntryIDs []int64
	resp, err := s.client.Do(req, &unstarredEntryIDs)
	if err != nil {
		return nil, resp, err
	}

	return unstarredEntryIDs, resp, nil
}
