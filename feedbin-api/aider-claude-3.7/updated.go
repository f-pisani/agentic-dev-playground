package feedbin

import (
	"net/http"
	"net/url"
	"time"
)

// GetUpdatedEntries retrieves updated entry IDs
func (c *Client) GetUpdatedEntries() ([]int64, error) {
	req, err := c.NewRequest(http.MethodGet, "/updated_entries.json", nil)
	if err != nil {
		return nil, err
	}
	
	var updatedIDs []int64
	_, err = c.Do(req, &updatedIDs)
	if err != nil {
		return nil, err
	}
	
	return updatedIDs, nil
}

// GetUpdatedEntriesSince retrieves updated entry IDs since a specific time
func (c *Client) GetUpdatedEntriesSince(since time.Time) ([]int64, error) {
	params := url.Values{}
	params.Set("since", FormatFeedbinTime(since))
	
	path := "/updated_entries.json?" + params.Encode()
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	
	var updatedIDs []int64
	_, err = c.Do(req, &updatedIDs)
	if err != nil {
		return nil, err
	}
	
	return updatedIDs, nil
}

// GetUpdatedEntriesContent retrieves the full content of updated entries
func (c *Client) GetUpdatedEntriesContent() ([]Entry, error) {
	// First get all updated entry IDs
	updatedIDs, err := c.GetUpdatedEntries()
	if err != nil {
		return nil, err
	}
	
	if len(updatedIDs) == 0 {
		return []Entry{}, nil
	}
	
	// Then get the entries
	return c.GetEntriesByIDs(updatedIDs)
}
