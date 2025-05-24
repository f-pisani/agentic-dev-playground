package feedbin

import (
	"fmt"
)

// GetFeed retrieves feed information by ID
func (c *Client) GetFeed(id int) (*Feed, error) {
	var feed Feed
	path := fmt.Sprintf("/feeds/%d.json", id)
	err := c.get(path, nil, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

// GetFeedEntries retrieves entries for a specific feed
func (c *Client) GetFeedEntries(feedID int, options ...QueryOption) ([]Entry, error) {
	var entries []Entry
	path := fmt.Sprintf("/feeds/%d/entries.json", feedID)
	query := buildQuery(options...)
	err := c.get(path, query, &entries)
	return entries, err
}
