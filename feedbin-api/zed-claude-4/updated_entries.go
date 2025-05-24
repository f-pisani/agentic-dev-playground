package feedbin

import (
	"context"
)

// GetUpdatedEntries retrieves the list of updated entry IDs
func (c *Client) GetUpdatedEntries(ctx context.Context, opts *UpdatedEntriesOptions) ([]int, error) {
	path := "/updated_entries.json"
	if opts != nil {
		params := buildQueryParams(opts)
		path = addQueryParams(path, params)
	}

	var entryIDs []int
	_, err := c.makeRequest(ctx, "GET", path, nil, &entryIDs)
	if err != nil {
		return nil, err
	}

	return entryIDs, nil
}

// MarkUpdatedEntriesAsRead marks the specified updated entry IDs as read
func (c *Client) MarkUpdatedEntriesAsRead(ctx context.Context, entryIDs []int) ([]int, error) {
	if len(entryIDs) == 0 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "at least one entry ID is required",
		}
	}

	req := &UpdatedEntriesRequest{
		UpdatedEntries: entryIDs,
	}

	var resultIDs []int
	_, err := c.makeRequest(ctx, "DELETE", "/updated_entries.json", req, &resultIDs)
	if err != nil {
		return nil, err
	}

	return resultIDs, nil
}

// MarkUpdatedEntriesAsReadPOST marks updated entries as read using the POST alternative endpoint
func (c *Client) MarkUpdatedEntriesAsReadPOST(ctx context.Context, entryIDs []int) ([]int, error) {
	if len(entryIDs) == 0 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "at least one entry ID is required",
		}
	}

	req := &UpdatedEntriesRequest{
		UpdatedEntries: entryIDs,
	}

	var resultIDs []int
	_, err := c.makeRequest(ctx, "POST", "/updated_entries/delete.json", req, &resultIDs)
	if err != nil {
		return nil, err
	}

	return resultIDs, nil
}

// Convenience methods without context (use background context)

func (c *Client) GetUpdatedEntriesWithoutContext(opts *UpdatedEntriesOptions) ([]int, error) {
	return c.GetUpdatedEntries(context.Background(), opts)
}

func (c *Client) MarkUpdatedEntriesAsReadWithoutContext(entryIDs []int) ([]int, error) {
	return c.MarkUpdatedEntriesAsRead(context.Background(), entryIDs)
}

func (c *Client) MarkUpdatedEntriesAsReadPOSTWithoutContext(entryIDs []int) ([]int, error) {
	return c.MarkUpdatedEntriesAsReadPOST(context.Background(), entryIDs)
}
