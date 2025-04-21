package feedbin

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// GetEntries retrieves entries with optional parameters
func (c *Client) GetEntries(params url.Values) ([]Entry, error) {
	path := "/entries.json"
	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}
	
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	
	var entries []Entry
	resp, err := c.Do(req, &entries)
	if err != nil {
		return nil, err
	}
	
	// Parse dates
	for i := range entries {
		if !entries[i].Published.IsZero() {
			entries[i].Published = entries[i].Published.UTC()
		}
		if !entries[i].CreatedAt.IsZero() {
			entries[i].CreatedAt = entries[i].CreatedAt.UTC()
		}
	}
	
	return entries, nil
}

// GetEntry retrieves a specific entry by ID
func (c *Client) GetEntry(id int64) (*Entry, error) {
	path := fmt.Sprintf("/entries/%d.json", id)
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	
	entry := new(Entry)
	_, err = c.Do(req, entry)
	if err != nil {
		return nil, err
	}
	
	// Parse dates
	if !entry.Published.IsZero() {
		entry.Published = entry.Published.UTC()
	}
	if !entry.CreatedAt.IsZero() {
		entry.CreatedAt = entry.CreatedAt.UTC()
	}
	
	return entry, nil
}

// GetFeedEntries retrieves entries for a specific feed
func (c *Client) GetFeedEntries(feedID int64, params url.Values) ([]Entry, error) {
	path := fmt.Sprintf("/feeds/%d/entries.json", feedID)
	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}
	
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	
	var entries []Entry
	resp, err := c.Do(req, &entries)
	if err != nil {
		return nil, err
	}
	
	// Parse dates
	for i := range entries {
		if !entries[i].Published.IsZero() {
			entries[i].Published = entries[i].Published.UTC()
		}
		if !entries[i].CreatedAt.IsZero() {
			entries[i].CreatedAt = entries[i].CreatedAt.UTC()
		}
	}
	
	return entries, nil
}

// GetEntriesSince retrieves entries published since a specific time
func (c *Client) GetEntriesSince(since time.Time, params url.Values) ([]Entry, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("since", FormatFeedbinTime(since))
	
	return c.GetEntries(params)
}

// GetEntriesByIDs retrieves entries by their IDs
func (c *Client) GetEntriesByIDs(ids []int64) ([]Entry, error) {
	if len(ids) == 0 {
		return []Entry{}, nil
	}
	
	params := url.Values{}
	for _, id := range ids {
		params.Add("ids", strconv.FormatInt(id, 10))
	}
	
	return c.GetEntries(params)
}

// GetEntryCount returns the total number of entries
func (c *Client) GetEntryCount() (int, error) {
	req, err := c.NewRequest(http.MethodGet, "/entries.json", nil)
	if err != nil {
		return 0, err
	}
	
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("API error: %s", resp.Status)
	}
	
	return GetTotalCount(resp), nil
}

// GetPaginatedEntries retrieves entries with pagination support
func (c *Client) GetPaginatedEntries(params url.Values, page int) ([]Entry, map[string]string, error) {
	if params == nil {
		params = url.Values{}
	}
	
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	
	path := fmt.Sprintf("/entries.json?%s", params.Encode())
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	
	var entries []Entry
	resp, err := c.Do(req, &entries)
	if err != nil {
		return nil, nil, err
	}
	
	// Parse dates
	for i := range entries {
		if !entries[i].Published.IsZero() {
			entries[i].Published = entries[i].Published.UTC()
		}
		if !entries[i].CreatedAt.IsZero() {
			entries[i].CreatedAt = entries[i].CreatedAt.UTC()
		}
	}
	
	links := GetPaginationLinks(resp)
	
	return entries, links, nil
}
