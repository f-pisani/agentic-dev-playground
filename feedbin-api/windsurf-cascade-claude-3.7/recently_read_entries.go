package feedbin

import (
	"net/http"
)

// RecentlyReadEntriesService handles operations related to recently read entries
type RecentlyReadEntriesService struct {
	client *Client
}

// List returns all recently read entry IDs
func (s *RecentlyReadEntriesService) List() ([]int, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/recently_read_entries.json", nil)
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

// Create marks the specified entry IDs as recently read
// Returns the entry IDs that were successfully marked as recently read
func (s *RecentlyReadEntriesService) Create(entryIDs []int) ([]int, error) {
	params := &RecentlyReadEntryIDs{
		RecentlyReadEntries: entryIDs,
	}

	req, err := s.client.NewRequest(http.MethodPost, "/recently_read_entries.json", params)
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
