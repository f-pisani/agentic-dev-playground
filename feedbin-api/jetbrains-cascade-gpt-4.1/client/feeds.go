package client

import (
	"fmt"
	"net/http"

	"feedbin-api/jetbrains-cascade-gpt-4.1/models"
)

// GetFeed fetches a feed by ID.
func (c *Client) GetFeed(id int64) (*models.Feed, error) {
	path := fmt.Sprintf("/feeds/%d.json", id)
	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var feed models.Feed
	err = c.do(req, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}
