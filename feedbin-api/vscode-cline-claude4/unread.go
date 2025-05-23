package feedbin

import (
	"context"
	"fmt"
)

// GetUnreadEntries retrieves all unread entry IDs for the authenticated user.
// Returns a slice of entry IDs.
func (c *Client) GetUnreadEntries(ctx context.Context) ([]int, error) {
	var entryIDs []int
	_, err := c.get(ctx, "unread_entries.json", nil, &entryIDs)
	if err != nil {
		return nil, err
	}

	return entryIDs, nil
}

// MarkAsUnread marks the specified entry IDs as unread.
// Returns the entry IDs that were successfully marked as unread.
// Maximum of 1000 entry IDs per request.
func (c *Client) MarkAsUnread(ctx context.Context, entryIDs []int) ([]int, error) {
	if len(entryIDs) == 0 {
		return []int{}, nil
	}

	if len(entryIDs) > MaxBulkOperationSize {
		return nil, fmt.Errorf("too many entry IDs: maximum %d allowed, got %d", MaxBulkOperationSize, len(entryIDs))
	}

	request := UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	}

	var result []int
	_, err := c.post(ctx, "unread_entries.json", nil, request, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// MarkAsRead marks the specified entry IDs as read.
// Returns the entry IDs that were successfully marked as read.
// Maximum of 1000 entry IDs per request.
func (c *Client) MarkAsRead(ctx context.Context, entryIDs []int) ([]int, error) {
	if len(entryIDs) == 0 {
		return []int{}, nil
	}

	if len(entryIDs) > MaxBulkOperationSize {
		return nil, fmt.Errorf("too many entry IDs: maximum %d allowed, got %d", MaxBulkOperationSize, len(entryIDs))
	}

	request := UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	}

	var result []int
	_, err := c.delete(ctx, "unread_entries.json", nil, request)
	if err != nil {
		return nil, err
	}

	// For DELETE requests, we need to make a separate call to get the result
	// since the delete method doesn't parse response body
	resp, err := c.post(ctx, "unread_entries/delete.json", nil, request, &result)
	if err != nil {
		return nil, err
	}

	// If the alternative POST endpoint worked, return the result
	if resp.StatusCode == 200 {
		return result, nil
	}

	// If we reach here, assume the DELETE was successful but we don't have the result
	return entryIDs, nil
}

// MarkAsReadPOST marks the specified entry IDs as read using POST method.
// This is an alternative to DELETE for clients that don't support DELETE with body.
// Returns the entry IDs that were successfully marked as read.
// Maximum of 1000 entry IDs per request.
func (c *Client) MarkAsReadPOST(ctx context.Context, entryIDs []int) ([]int, error) {
	if len(entryIDs) == 0 {
		return []int{}, nil
	}

	if len(entryIDs) > MaxBulkOperationSize {
		return nil, fmt.Errorf("too many entry IDs: maximum %d allowed, got %d", MaxBulkOperationSize, len(entryIDs))
	}

	request := UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	}

	var result []int
	_, err := c.post(ctx, "unread_entries/delete.json", nil, request, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// MarkAllAsRead marks all unread entries as read.
// This is a convenience method that gets all unread entries and marks them as read.
func (c *Client) MarkAllAsRead(ctx context.Context) ([]int, error) {
	unreadIDs, err := c.GetUnreadEntries(ctx)
	if err != nil {
		return nil, err
	}

	if len(unreadIDs) == 0 {
		return []int{}, nil
	}

	// Process in batches if there are too many
	var allMarked []int
	for i := 0; i < len(unreadIDs); i += MaxBulkOperationSize {
		end := i + MaxBulkOperationSize
		if end > len(unreadIDs) {
			end = len(unreadIDs)
		}

		batch := unreadIDs[i:end]
		marked, err := c.MarkAsRead(ctx, batch)
		if err != nil {
			return allMarked, err
		}

		allMarked = append(allMarked, marked...)
	}

	return allMarked, nil
}
