package client

import (
	"net/http"

	"feedbin-api/jetbrains-cascade-gpt-4.1/models"
)

// GetUnreadEntries fetches all unread entry IDs for the user.
func (c *Client) GetUnreadEntries() ([]models.UnreadEntry, error) {
	req, err := c.newRequest(http.MethodGet, "/unread_entries.json", nil)
	if err != nil {
		return nil, err
	}
	var entries models.UnreadEntriesResponse
	err = c.do(req, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
