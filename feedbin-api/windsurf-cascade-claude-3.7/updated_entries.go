package feedbin

import (
	"net/http"
	"net/url"
	"time"
)

// UpdatedEntriesService handles operations related to updated entries
type UpdatedEntriesService struct {
	client *Client
}

// List returns all updated entry IDs
// If since is provided, returns only entries updated after the specified time
func (s *UpdatedEntriesService) List(since *time.Time) ([]int, error) {
	path := "/updated_entries.json"

	// Add since parameter if provided
	if since != nil {
		values := url.Values{}
		values.Add("since", since.Format(time.RFC3339Nano))
		path = path + "?" + values.Encode()
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
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

// Delete marks the specified updated entry IDs as read
// Returns the entry IDs that were successfully marked as read
func (s *UpdatedEntriesService) Delete(entryIDs []int) ([]int, error) {
	params := &UpdatedEntryIDs{
		UpdatedEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodDelete, "/updated_entries.json", params)
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
func (s *UpdatedEntriesService) DeleteWithPost(entryIDs []int) ([]int, error) {
	params := &UpdatedEntryIDs{
		UpdatedEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodPost, "/updated_entries/delete.json", params)
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
