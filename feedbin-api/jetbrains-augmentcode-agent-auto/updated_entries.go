package feedbin

import (
	"net/http"
)

// UpdatedEntriesService handles communication with the updated entries related
// methods of the Feedbin API.
type UpdatedEntriesService struct {
	client *Client
}

// List returns all updated entry IDs.
func (s *UpdatedEntriesService) List() ([]int64, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/v2/updated_entries.json", nil)
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

// UpdatedEntriesRequest represents a request to mark entries as read.
type UpdatedEntriesRequest struct {
	UpdatedEntries []int64 `json:"updated_entries"`
}

// Delete marks updated entries as read.
func (s *UpdatedEntriesService) Delete(entryIDs []int64) ([]int64, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodDelete, "/v2/updated_entries.json", &UpdatedEntriesRequest{
		UpdatedEntries: entryIDs,
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

// DeleteAlternative marks updated entries as read using the POST alternative for DELETE.
func (s *UpdatedEntriesService) DeleteAlternative(entryIDs []int64) ([]int64, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/v2/updated_entries/delete.json", &UpdatedEntriesRequest{
		UpdatedEntries: entryIDs,
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
