package client

import (
	"fmt"
	"net/url"
	"time"
)

// RecentlyReadService handles communication with the recently read entries
// methods of the Feedbin API
type RecentlyReadService struct {
	client *Client
}

// RecentlyReadListOptions specifies the optional parameters to the
// RecentlyReadService.List method
type RecentlyReadListOptions struct {
	Since time.Time // Get entries created after this time
}

// List returns all recently read entry IDs
func (s *RecentlyReadService) List(opts *RecentlyReadListOptions) ([]int64, error) {
	u := "recently_read_entries.json"
	u, err := addRecentlyReadOptions(u, opts)
	if err != nil {
		return nil, err
	}

	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var entryIDs []int64
	_, err = s.client.do(req, &entryIDs)
	if err != nil {
		return nil, err
	}

	return entryIDs, nil
}

// Create adds entries to the recently read list
func (s *RecentlyReadService) Create(entryIDs []int64) ([]int64, error) {
	if len(entryIDs) > 1000 {
		return nil, fmt.Errorf("maximum of 1000 entry IDs can be added to recently read in a single request")
	}

	type recentlyReadRequest struct {
		RecentlyReadEntries []int64 `json:"recently_read_entries"`
	}

	req, err := s.client.newRequest("POST", "recently_read_entries.json", &recentlyReadRequest{
		RecentlyReadEntries: entryIDs,
	})
	if err != nil {
		return nil, err
	}

	var processedIDs []int64
	_, err = s.client.do(req, &processedIDs)
	if err != nil {
		return nil, err
	}

	return processedIDs, nil
}

// addRecentlyReadOptions adds the parameters in opts as URL query parameters to s.
// opts is a pointer to a RecentlyReadListOptions struct.
func addRecentlyReadOptions(s string, opts *RecentlyReadListOptions) (string, error) {
	if opts == nil {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	q := u.Query()
	if !opts.Since.IsZero() {
		q.Set("since", opts.Since.Format(time.RFC3339Nano))
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
