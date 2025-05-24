package feedbin

import (
	"context"
)

// GetRecentlyReadEntries retrieves the list of recently read entry IDs
func (c *Client) GetRecentlyReadEntries(ctx context.Context) ([]int, error) {
	var entryIDs []int
	_, err := c.makeRequest(ctx, "GET", "/recently_read_entries.json", nil, &entryIDs)
	if err != nil {
		return nil, err
	}

	return entryIDs, nil
}

// CreateRecentlyReadEntries marks the specified entry IDs as recently read
func (c *Client) CreateRecentlyReadEntries(ctx context.Context, entryIDs []int) ([]int, error) {
	if len(entryIDs) == 0 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "at least one entry ID is required",
		}
	}

	req := &RecentlyReadEntriesRequest{
		RecentlyReadEntries: entryIDs,
	}

	var resultIDs []int
	_, err := c.makeRequest(ctx, "POST", "/recently_read_entries.json", req, &resultIDs)
	if err != nil {
		return nil, err
	}

	return resultIDs, nil
}

// Convenience methods without context (use background context)

func (c *Client) GetRecentlyReadEntriesWithoutContext() ([]int, error) {
	return c.GetRecentlyReadEntries(context.Background())
}

func (c *Client) CreateRecentlyReadEntriesWithoutContext(entryIDs []int) ([]int, error) {
	return c.CreateRecentlyReadEntries(context.Background(), entryIDs)
}
