package feedbin

import (
	"context"
)

// GetRecentlyReadEntries retrieves all recently read entry IDs for the authenticated user.
// Returns a slice of entry IDs.
func (c *Client) GetRecentlyReadEntries(ctx context.Context) ([]int, error) {
	var entryIDs []int
	_, err := c.get(ctx, "recently_read_entries.json", nil, &entryIDs)
	if err != nil {
		return nil, err
	}

	return entryIDs, nil
}
