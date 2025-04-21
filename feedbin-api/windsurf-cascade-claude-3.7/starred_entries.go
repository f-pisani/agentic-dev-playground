package feedbin

import (
	"net/http"
)

// StarredEntriesService handles operations related to starred entries
type StarredEntriesService struct {
	client *Client
}

// List returns all starred entry IDs
func (s *StarredEntriesService) List() ([]int, error) {
	req, err := s.client.NewRequest(http.MethodGet, "starred_entries.json", nil)
	if err != nil {
		return nil, err
	}

	var entryIDs []int
	_, err = s.client.Do(req, &entryIDs)
	if err != nil {
		return nil, err
	}

	return entryIDs, nil
}

// Create stars the specified entry IDs
// Returns the entry IDs that were successfully starred
func (s *StarredEntriesService) Create(entryIDs []int) ([]int, error) {
	params := &StarredEntryIDs{
		StarredEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodPost, "starred_entries.json", params)
	if err != nil {
		return nil, err
	}

	var responseIDs []int
	_, err = s.client.Do(req, &responseIDs)
	if err != nil {
		return nil, err
	}

	return responseIDs, nil
}

// Delete unstars the specified entry IDs
// Returns the entry IDs that were successfully unstarred
func (s *StarredEntriesService) Delete(entryIDs []int) ([]int, error) {
	params := &StarredEntryIDs{
		StarredEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodDelete, "starred_entries.json", params)
	if err != nil {
		return nil, err
	}

	var responseIDs []int
	_, err = s.client.Do(req, &responseIDs)
	if err != nil {
		return nil, err
	}

	return responseIDs, nil
}

// DeleteWithPost is an alternative to Delete that uses POST instead of DELETE
// This is useful for clients that don't support DELETE requests with a body
// Returns the entry IDs that were successfully unstarred
func (s *StarredEntriesService) DeleteWithPost(entryIDs []int) ([]int, error) {
	params := &StarredEntryIDs{
		StarredEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodPost, "starred_entries/delete.json", params)
	if err != nil {
		return nil, err
	}

	var responseIDs []int
	_, err = s.client.Do(req, &responseIDs)
	if err != nil {
		return nil, err
	}

	return responseIDs, nil
}
