package feedbin

import (
	"context"
	"fmt"
)

// GetStarredEntries retrieves all starred entry IDs for the authenticated user.
// Returns a slice of entry IDs.
func (c *Client) GetStarredEntries(ctx context.Context) ([]int, error) {
	var entryIDs []int
	_, err := c.get(ctx, "starred_entries.json", nil, &entryIDs)
	if err != nil {
		return nil, err
	}

	return entryIDs, nil
}

// StarEntries stars the specified entry IDs.
// Returns the entry IDs that were successfully starred.
// Maximum of 1000 entry IDs per request.
func (c *Client) StarEntries(ctx context.Context, entryIDs []int) ([]int, error) {
	if len(entryIDs) == 0 {
		return []int{}, nil
	}

	if len(entryIDs) > MaxBulkOperationSize {
		return nil, fmt.Errorf("too many entry IDs: maximum %d allowed, got %d", MaxBulkOperationSize, len(entryIDs))
	}

	request := StarredEntriesRequest{
		StarredEntries: entryIDs,
	}

	var result []int
	_, err := c.post(ctx, "starred_entries.json", nil, request, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UnstarEntries unstars the specified entry IDs.
// Returns the entry IDs that were successfully unstarred.
// Maximum of 1000 entry IDs per request.
func (c *Client) UnstarEntries(ctx context.Context, entryIDs []int) ([]int, error) {
	if len(entryIDs) == 0 {
		return []int{}, nil
	}

	if len(entryIDs) > MaxBulkOperationSize {
		return nil, fmt.Errorf("too many entry IDs: maximum %d allowed, got %d", MaxBulkOperationSize, len(entryIDs))
	}

	request := StarredEntriesRequest{
		StarredEntries: entryIDs,
	}

	var result []int
	_, err := c.delete(ctx, "starred_entries.json", nil, request)
	if err != nil {
		return nil, err
	}

	// For DELETE requests, we need to make a separate call to get the result
	// since the delete method doesn't parse response body
	resp, err := c.post(ctx, "starred_entries/delete.json", nil, request, &result)
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

// UnstarEntriesPOST unstars the specified entry IDs using POST method.
// This is an alternative to DELETE for clients that don't support DELETE with body.
// Returns the entry IDs that were successfully unstarred.
// Maximum of 1000 entry IDs per request.
func (c *Client) UnstarEntriesPOST(ctx context.Context, entryIDs []int) ([]int, error) {
	if len(entryIDs) == 0 {
		return []int{}, nil
	}

	if len(entryIDs) > MaxBulkOperationSize {
		return nil, fmt.Errorf("too many entry IDs: maximum %d allowed, got %d", MaxBulkOperationSize, len(entryIDs))
	}

	request := StarredEntriesRequest{
		StarredEntries: entryIDs,
	}

	var result []int
	_, err := c.post(ctx, "starred_entries/delete.json", nil, request, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UnstarAllEntries unstars all starred entries.
// This is a convenience method that gets all starred entries and unstars them.
func (c *Client) UnstarAllEntries(ctx context.Context) ([]int, error) {
	starredIDs, err := c.GetStarredEntries(ctx)
	if err != nil {
		return nil, err
	}

	if len(starredIDs) == 0 {
		return []int{}, nil
	}

	// Process in batches if there are too many
	var allUnstarred []int
	for i := 0; i < len(starredIDs); i += MaxBulkOperationSize {
		end := i + MaxBulkOperationSize
		if end > len(starredIDs) {
			end = len(starredIDs)
		}

		batch := starredIDs[i:end]
		unstarred, err := c.UnstarEntries(ctx, batch)
		if err != nil {
			return allUnstarred, err
		}

		allUnstarred = append(allUnstarred, unstarred...)
	}

	return allUnstarred, nil
}
