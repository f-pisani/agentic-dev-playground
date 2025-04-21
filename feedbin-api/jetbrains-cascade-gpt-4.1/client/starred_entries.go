package client

import (
	"net/http"

	"feedbin-api/jetbrains-cascade-gpt-4.1/models"
)

// GetStarredEntries fetches all starred entry IDs for the user.
func (c *Client) GetStarredEntries() ([]models.StarredEntry, error) {
	req, err := c.newRequest(http.MethodGet, "/starred_entries.json", nil)
	if err != nil {
		return nil, err
	}
	var entries models.StarredEntriesResponse
	err = c.do(req, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
