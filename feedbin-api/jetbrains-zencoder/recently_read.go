package feedbin

import (
	"net/http"
	"net/url"
	"time"
)

// GetRecentlyReadEntries retrieves the IDs of recently read entries.
func (c *Client) GetRecentlyReadEntries(since *time.Time) ([]int, error) {
	// Build query parameters
	query := url.Values{}

	if since != nil {
		query.Set("since", since.Format(time.RFC3339Nano))
	}

	// Make the request
	resp, err := c.doRequest(http.MethodGet, "/recently_read_entries.json", nil, query)
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

// MarkEntriesAsRecentlyRead marks the specified entries as recently read.
func (c *Client) MarkEntriesAsRecentlyRead(entryIDs []int) error {
	// Create the request body
	reqBody := EntryIDs{
		IDs: entryIDs,
	}

	// Make the request
	resp, err := c.doRequest(http.MethodPost, "/recently_read_entries.json", reqBody, nil)
	if err != nil {
		return err
	}

	// Check for success (200 OK)
	if resp.StatusCode != http.StatusOK {
		return parseResponse(resp, nil)
	}

	return nil
}
