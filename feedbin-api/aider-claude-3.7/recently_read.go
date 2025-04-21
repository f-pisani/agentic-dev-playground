package feedbin

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// GetRecentlyReadEntries retrieves recently read entry IDs
func (c *Client) GetRecentlyReadEntries() ([]int64, error) {
	req, err := c.NewRequest(http.MethodGet, "/recently_read_entries.json", nil)
	if err != nil {
		return nil, err
	}
	
	var recentlyReadIDs []int64
	_, err = c.Do(req, &recentlyReadIDs)
	if err != nil {
		return nil, err
	}
	
	return recentlyReadIDs, nil
}

// GetRecentlyReadEntriesSince retrieves recently read entry IDs since a specific time
func (c *Client) GetRecentlyReadEntriesSince(since time.Time) ([]int64, error) {
	params := url.Values{}
	params.Set("since", FormatFeedbinTime(since))
	
	path := "/recently_read_entries.json?" + params.Encode()
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	
	var recentlyReadIDs []int64
	_, err = c.Do(req, &recentlyReadIDs)
	if err != nil {
		return nil, err
	}
	
	return recentlyReadIDs, nil
}

// MarkEntriesAsRecentlyRead marks entries as recently read
func (c *Client) MarkEntriesAsRecentlyRead(entryIDs []int64) error {
	if len(entryIDs) == 0 {
		return nil
	}
	
	recentlyReadReq := &RecentlyReadEntryRequest{
		EntryIDs: entryIDs,
	}
	
	req, err := c.NewRequest(http.MethodPost, "/recently_read_entries.json", recentlyReadReq)
	if err != nil {
		return err
	}
	
	_, err = c.Do(req, nil)
	return err
}

// GetRecentlyReadEntriesContent retrieves the full content of recently read entries
func (c *Client) GetRecentlyReadEntriesContent() ([]Entry, error) {
	// First get all recently read entry IDs
	recentlyReadIDs, err := c.GetRecentlyReadEntries()
	if err != nil {
		return nil, err
	}
	
	if len(recentlyReadIDs) == 0 {
		return []Entry{}, nil
	}
	
	// Then get the entries
	return c.GetEntriesByIDs(recentlyReadIDs)
}
