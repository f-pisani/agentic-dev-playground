// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"net/http"
)

// StarredService provides access to the starred entries endpoints of the Feedbin API.
type StarredService struct {
	client *Client
}

// StarredEntriesRequest represents a request to star or unstar entries.
type StarredEntriesRequest struct {
	StarredEntries []int `json:"starred_entries"`
}

// GetStarredEntries retrieves all starred entry IDs.
func (s *StarredService) GetStarredEntries() ([]int, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "/starred_entries.json", nil)
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

// StarEntries stars the specified entries.
func (s *StarredService) StarEntries(entryIDs []int) ([]int, *http.Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, nil
	}

	if len(entryIDs) > 1000 {
		return nil, nil, &APIError{
			Message: "maximum of 1,000 entry IDs allowed",
		}
	}

	req, err := s.client.NewRequest("POST", "/starred_entries.json", &StarredEntriesRequest{
		StarredEntries: entryIDs,
	})
	if err != nil {
		return nil, nil, err
	}

	var starredEntryIDs []int
	resp, err := s.client.Do(req, &starredEntryIDs)
	if err != nil {
		return nil, resp, err
	}

	return starredEntryIDs, resp, nil
}

// UnstarEntries unstars the specified entries.
func (s *StarredService) UnstarEntries(entryIDs []int) ([]int, *http.Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, nil
	}

	if len(entryIDs) > 1000 {
		return nil, nil, &APIError{
			Message: "maximum of 1,000 entry IDs allowed",
		}
	}

	req, err := s.client.NewRequest("DELETE", "/starred_entries.json", &StarredEntriesRequest{
		StarredEntries: entryIDs,
	})
	if err != nil {
		return nil, nil, err
	}

	var unstarredEntryIDs []int
	resp, err := s.client.Do(req, &unstarredEntryIDs)
	if err != nil {
		return nil, resp, err
	}

	return unstarredEntryIDs, resp, nil
}

// UnstarEntriesAlternative unstars the specified entries using the POST alternative.
// Some clients like Android don't easily allow a body with a DELETE request.
func (s *StarredService) UnstarEntriesAlternative(entryIDs []int) ([]int, *http.Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, nil
	}

	if len(entryIDs) > 1000 {
		return nil, nil, &APIError{
			Message: "maximum of 1,000 entry IDs allowed",
		}
	}

	req, err := s.client.NewRequest("POST", "/starred_entries/delete.json", &StarredEntriesRequest{
		StarredEntries: entryIDs,
	})
	if err != nil {
		return nil, nil, err
	}

	var unstarredEntryIDs []int
	resp, err := s.client.Do(req, &unstarredEntryIDs)
	if err != nil {
		return nil, resp, err
	}

	return unstarredEntryIDs, resp, nil
}
