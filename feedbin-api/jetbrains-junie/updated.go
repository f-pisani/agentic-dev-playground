// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"net/url"
	"strconv"
	"time"
)

// UpdatedService handles communication with the updated entries related
// methods of the Feedbin API.
type UpdatedService struct {
	client *Client
}

// GetUpdatedEntries retrieves updated entries.
func (s *UpdatedService) GetUpdatedEntries(since time.Time) ([]int, error) {
	path := "/updated_entries.json"

	// Add query parameters
	params := url.Values{}
	params.Add("since", since.Format(time.RFC3339Nano))

	// Append query parameters to path
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", path, nil)
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

// GetUpdatedEntriesByPage retrieves updated entries with pagination.
func (s *UpdatedService) GetUpdatedEntriesByPage(since time.Time, page int) ([]int, *PaginationLinks, error) {
	path := "/updated_entries.json"

	// Add query parameters
	params := url.Values{}
	params.Add("since", since.Format(time.RFC3339Nano))
	params.Add("page", strconv.Itoa(page))

	// Append query parameters to path
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var entryIDs []int
	resp, err := s.client.Do(req, &entryIDs)
	if err != nil {
		return nil, nil, err
	}

	// Parse pagination links
	links := ParseLinkHeader(resp)

	return entryIDs, links, nil
}
