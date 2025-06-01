package feedbin

import (
	"net/http"
	"net/url"
	"time"
)

// GetUpdatedEntries retrieves the IDs of entries that have been updated.
func (c *Client) GetUpdatedEntries(since *time.Time) ([]int, error) {
	// Build query parameters
	query := url.Values{}

	if since != nil {
		query.Set("since", since.Format(time.RFC3339Nano))
	}

	// Make the request
	resp, err := c.doRequest(http.MethodGet, "/updated_entries.json", nil, query)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var ids []int
	if err := parseResponse(resp, &ids); err != nil {
		return nil, err
	}

	return ids, nil
}
