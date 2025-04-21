package feedbin

import (
	"net/http"
)

// GetStarredEntries retrieves all starred entry IDs
func (c *Client) GetStarredEntries() ([]int64, error) {
	req, err := c.NewRequest(http.MethodGet, "/v2/starred_entries.json", nil)
	if err != nil {
		return nil, err
	}
	
	var starredIDs []int64
	_, err = c.Do(req, &starredIDs)
	if err != nil {
		return nil, err
	}
	
	return starredIDs, nil
}

// StarEntries marks entries as starred
func (c *Client) StarEntries(entryIDs []int64) error {
	if len(entryIDs) == 0 {
		return nil
	}
	
	starredReq := &StarredEntryRequest{
		StarredEntries: entryIDs,
	}
	
	req, err := c.NewRequest(http.MethodPost, "/v2/starred_entries.json", starredReq)
	if err != nil {
		return err
	}
	
	_, err = c.Do(req, nil)
	return err
}

// UnstarEntries removes the star from entries
func (c *Client) UnstarEntries(entryIDs []int64) error {
	if len(entryIDs) == 0 {
		return nil
	}
	
	starredReq := &StarredEntryRequest{
		StarredEntries: entryIDs,
	}
	
	req, err := c.NewRequest(http.MethodDelete, "/v2/starred_entries.json", starredReq)
	if err != nil {
		return err
	}
	
	_, err = c.Do(req, nil)
	return err
}

// GetStarredCount returns the total number of starred entries
func (c *Client) GetStarredCount() (int, error) {
	starredIDs, err := c.GetStarredEntries()
	if err != nil {
		return 0, err
	}
	
	return len(starredIDs), nil
}

// GetStarredEntriesContent retrieves the full content of all starred entries
func (c *Client) GetStarredEntriesContent() ([]Entry, error) {
	// First get all starred entry IDs
	starredIDs, err := c.GetStarredEntries()
	if err != nil {
		return nil, err
	}
	
	if len(starredIDs) == 0 {
		return []Entry{}, nil
	}
	
	// Then get the entries
	return c.GetEntriesByIDs(starredIDs)
}
