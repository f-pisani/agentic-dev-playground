// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"fmt"
	"net/http"
)

// RecentlyReadService handles communication with the recently read entries related
// methods of the Feedbin API.
type RecentlyReadService struct {
	client *Client
}

// GetRecentlyRead retrieves recently read entries.
func (s *RecentlyReadService) GetRecentlyRead() ([]int, error) {
	req, err := s.client.NewRequest("GET", "/recently_read_entries.json", nil)
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

// AddRecentlyRead adds entries to the recently read list.
func (s *RecentlyReadService) AddRecentlyRead(entryIDs []int) error {
	body := map[string]interface{}{
		"recently_read_entries": entryIDs,
	}

	req, err := s.client.NewRequest("POST", "/recently_read_entries.json", body)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// DeleteRecentlyRead removes entries from the recently read list.
func (s *RecentlyReadService) DeleteRecentlyRead(entryIDs []int) error {
	body := map[string]interface{}{
		"recently_read_entries": entryIDs,
	}

	req, err := s.client.NewRequest("DELETE", "/recently_read_entries.json", body)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
