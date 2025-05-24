package feedbin

import (
	"context"
	"fmt"
)

// GetFeed retrieves a specific feed by ID
func (c *Client) GetFeed(ctx context.Context, id int) (*Feed, error) {
	path := fmt.Sprintf("/feeds/%d.json", id)

	var feed Feed
	_, err := c.makeRequest(ctx, "GET", path, nil, &feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}

// GetFeedWithoutContext retrieves a feed without context (uses background context)
func (c *Client) GetFeedWithoutContext(id int) (*Feed, error) {
	return c.GetFeed(context.Background(), id)
}
