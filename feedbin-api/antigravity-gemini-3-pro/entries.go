package feedbin

import (
	"context"
	"fmt"
)

// GetEntries fetches entries for the user.
func (c *Client) GetEntries(ctx context.Context, opt *PaginationOptions) ([]Entry, *Response, error) {
	u, err := addOptions("entries.json", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var entries []Entry
	resp, err := c.Do(ctx, req, &entries)
	if err != nil {
		return nil, resp, err
	}

	return entries, resp, nil
}

// GetEntry fetches a single entry by ID.
func (c *Client) GetEntry(ctx context.Context, id int64) (*Entry, *Response, error) {
	u := fmt.Sprintf("entries/%d.json", id)
	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var entry Entry
	resp, err := c.Do(ctx, req, &entry)
	if err != nil {
		return nil, resp, err
	}

	return &entry, resp, nil
}

// GetEntriesForFeed fetches entries for a specific feed.
func (c *Client) GetEntriesForFeed(ctx context.Context, feedID int64, opt *PaginationOptions) ([]Entry, *Response, error) {
	u := fmt.Sprintf("feeds/%d/entries.json", feedID)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var entries []Entry
	resp, err := c.Do(ctx, req, &entries)
	if err != nil {
		return nil, resp, err
	}

	return entries, resp, nil
}
