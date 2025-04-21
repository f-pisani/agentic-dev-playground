package client

import (
	"net/http"

	"feedbin-api/jetbrains-cascade-gpt-4.1/models"
)

// GetRecentlyReadEntries fetches all recently read entry IDs for the user.
func (c *Client) GetRecentlyReadEntries() ([]models.RecentlyReadEntry, error) {
	req, err := c.newRequest(http.MethodGet, "/recently_read_entries.json", nil)
	if err != nil {
		return nil, err
	}
	var entries models.RecentlyReadEntriesResponse
	err = c.do(req, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
