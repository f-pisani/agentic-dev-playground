package feedbin

import (
	"context"
	"net/url"
	"time"
)

// GetUpdatedEntries retrieves all updated entry IDs for the authenticated user.
// Returns a slice of entry IDs.
func (c *Client) GetUpdatedEntries(ctx context.Context) ([]int, error) {
	var entryIDs []int
	_, err := c.get(ctx, "updated_entries.json", nil, &entryIDs)
	if err != nil {
		return nil, err
	}

	return entryIDs, nil
}

// GetUpdatedEntriesSince retrieves updated entry IDs since a specific time.
// Returns a slice of entry IDs.
func (c *Client) GetUpdatedEntriesSince(ctx context.Context, since time.Time) ([]int, error) {
	params := url.Values{}
	params.Set("since", since.Format("2006-01-02T15:04:05.000000Z"))

	var entryIDs []int
	_, err := c.get(ctx, "updated_entries.json", params, &entryIDs)
	if err != nil {
		return nil, err
	}

	return entryIDs, nil
}
