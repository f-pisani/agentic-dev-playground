package feedbin

import (
	"net/http"
)

// UnreadEntriesService handles operations related to unread entries
type UnreadEntriesService struct {
	client *Client
}

// List returns all unread entry IDs
func (s *UnreadEntriesService) List() ([]int, error) {
	req, err := s.client.NewRequest(http.MethodGet, "unread_entries.json", nil)
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

// Create marks the specified entry IDs as unread
// Returns the entry IDs that were successfully marked as unread
func (s *UnreadEntriesService) Create(entryIDs []int) ([]int, error) {
	params := &UnreadEntryIDs{
		UnreadEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodPost, "unread_entries.json", params)
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

// Delete marks the specified entry IDs as read
// Returns the entry IDs that were successfully marked as read
func (s *UnreadEntriesService) Delete(entryIDs []int) ([]int, error) {
	params := &UnreadEntryIDs{
		UnreadEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodDelete, "unread_entries.json", params)
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
// Returns the entry IDs that were successfully marked as read
func (s *UnreadEntriesService) DeleteWithPost(entryIDs []int) ([]int, error) {
	params := &UnreadEntryIDs{
		UnreadEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodPost, "unread_entries/delete.json", params)
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
