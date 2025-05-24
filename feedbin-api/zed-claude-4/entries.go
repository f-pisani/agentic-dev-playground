package feedbin

import (
	"context"
	"fmt"
)

// GetEntries retrieves entries with optional filtering and pagination
func (c *Client) GetEntries(ctx context.Context, opts *EntryOptions) ([]Entry, *PaginationInfo, error) {
	path := "/entries.json"
	if opts != nil {
		params := buildQueryParams(opts)
		path = addQueryParams(path, params)
	}

	var entries []Entry
	resp, err := c.makeRequest(ctx, "GET", path, nil, &entries)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	pagination := extractPaginationInfo(resp)
	return entries, pagination, nil
}

// GetFeedEntries retrieves entries for a specific feed
func (c *Client) GetFeedEntries(ctx context.Context, feedID int, opts *EntryOptions) ([]Entry, *PaginationInfo, error) {
	path := fmt.Sprintf("/feeds/%d/entries.json", feedID)
	if opts != nil {
		params := buildQueryParams(opts)
		path = addQueryParams(path, params)
	}

	var entries []Entry
	resp, err := c.makeRequest(ctx, "GET", path, nil, &entries)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	pagination := extractPaginationInfo(resp)
	return entries, pagination, nil
}

// GetEntry retrieves a specific entry by ID
func (c *Client) GetEntry(ctx context.Context, id int, opts *EntryOptions) (*Entry, error) {
	path := fmt.Sprintf("/entries/%d.json", id)
	if opts != nil {
		params := buildQueryParams(opts)
		path = addQueryParams(path, params)
	}

	var entry Entry
	_, err := c.makeRequest(ctx, "GET", path, nil, &entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// GetEntriesFromURL retrieves entries from a full URL (used for pagination)
func (c *Client) GetEntriesFromURL(ctx context.Context, url string) ([]Entry, *PaginationInfo, error) {
	var entries []Entry
	resp, err := c.GetFromURL(ctx, url, &entries)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	pagination := extractPaginationInfo(resp)
	return entries, pagination, nil
}

// GetUnreadEntries retrieves the list of unread entry IDs
func (c *Client) GetUnreadEntries(ctx context.Context) ([]int, error) {
	var entryIDs []int
	_, err := c.makeRequest(ctx, "GET", "/unread_entries.json", nil, &entryIDs)
	if err != nil {
		return nil, err
	}

	return entryIDs, nil
}

// MarkAsUnread marks the specified entry IDs as unread
func (c *Client) MarkAsUnread(ctx context.Context, entryIDs []int) ([]int, error) {
	if len(entryIDs) == 0 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "at least one entry ID is required",
		}
	}

	if len(entryIDs) > 1000 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "maximum of 1,000 entry IDs allowed per request",
		}
	}

	req := &UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	}

	var resultIDs []int
	_, err := c.makeRequest(ctx, "POST", "/unread_entries.json", req, &resultIDs)
	if err != nil {
		return nil, err
	}

	return resultIDs, nil
}

// MarkAsRead marks the specified entry IDs as read
func (c *Client) MarkAsRead(ctx context.Context, entryIDs []int) ([]int, error) {
	if len(entryIDs) == 0 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "at least one entry ID is required",
		}
	}

	if len(entryIDs) > 1000 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "maximum of 1,000 entry IDs allowed per request",
		}
	}

	req := &UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	}

	var resultIDs []int
	_, err := c.makeRequest(ctx, "DELETE", "/unread_entries.json", req, &resultIDs)
	if err != nil {
		return nil, err
	}

	return resultIDs, nil
}

// MarkAsReadPOST marks entries as read using the POST alternative endpoint
func (c *Client) MarkAsReadPOST(ctx context.Context, entryIDs []int) ([]int, error) {
	if len(entryIDs) == 0 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "at least one entry ID is required",
		}
	}

	if len(entryIDs) > 1000 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "maximum of 1,000 entry IDs allowed per request",
		}
	}

	req := &UnreadEntriesRequest{
		UnreadEntries: entryIDs,
	}

	var resultIDs []int
	_, err := c.makeRequest(ctx, "POST", "/unread_entries/delete.json", req, &resultIDs)
	if err != nil {
		return nil, err
	}

	return resultIDs, nil
}

// GetStarredEntries retrieves the list of starred entry IDs
func (c *Client) GetStarredEntries(ctx context.Context) ([]int, error) {
	var entryIDs []int
	_, err := c.makeRequest(ctx, "GET", "/starred_entries.json", nil, &entryIDs)
	if err != nil {
		return nil, err
	}

	return entryIDs, nil
}

// StarEntries stars the specified entry IDs
func (c *Client) StarEntries(ctx context.Context, entryIDs []int) ([]int, error) {
	if len(entryIDs) == 0 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "at least one entry ID is required",
		}
	}

	if len(entryIDs) > 1000 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "maximum of 1,000 entry IDs allowed per request",
		}
	}

	req := &StarredEntriesRequest{
		StarredEntries: entryIDs,
	}

	var resultIDs []int
	_, err := c.makeRequest(ctx, "POST", "/starred_entries.json", req, &resultIDs)
	if err != nil {
		return nil, err
	}

	return resultIDs, nil
}

// UnstarEntries unstars the specified entry IDs
func (c *Client) UnstarEntries(ctx context.Context, entryIDs []int) ([]int, error) {
	if len(entryIDs) == 0 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "at least one entry ID is required",
		}
	}

	if len(entryIDs) > 1000 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "maximum of 1,000 entry IDs allowed per request",
		}
	}

	req := &StarredEntriesRequest{
		StarredEntries: entryIDs,
	}

	var resultIDs []int
	_, err := c.makeRequest(ctx, "DELETE", "/starred_entries.json", req, &resultIDs)
	if err != nil {
		return nil, err
	}

	return resultIDs, nil
}

// UnstarEntriesPOST unstars entries using the POST alternative endpoint
func (c *Client) UnstarEntriesPOST(ctx context.Context, entryIDs []int) ([]int, error) {
	if len(entryIDs) == 0 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "at least one entry ID is required",
		}
	}

	if len(entryIDs) > 1000 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "maximum of 1,000 entry IDs allowed per request",
		}
	}

	req := &StarredEntriesRequest{
		StarredEntries: entryIDs,
	}

	var resultIDs []int
	_, err := c.makeRequest(ctx, "POST", "/starred_entries/delete.json", req, &resultIDs)
	if err != nil {
		return nil, err
	}

	return resultIDs, nil
}

// GetEntriesByIDs is a convenience method to get multiple entries by their IDs
func (c *Client) GetEntriesByIDs(ctx context.Context, entryIDs []int, opts *EntryOptions) ([]Entry, error) {
	if len(entryIDs) == 0 {
		return []Entry{}, nil
	}

	if len(entryIDs) > 100 {
		return nil, &ValidationError{
			Field:   "entry_ids",
			Message: "maximum of 100 entry IDs allowed per request",
		}
	}

	// Create new options or copy existing ones
	newOpts := &EntryOptions{}
	if opts != nil {
		*newOpts = *opts
	}
	newOpts.IDs = entryIDs

	entries, _, err := c.GetEntries(ctx, newOpts)
	return entries, err
}

// Convenience methods without context (use background context)

func (c *Client) GetEntriesWithoutContext(opts *EntryOptions) ([]Entry, *PaginationInfo, error) {
	return c.GetEntries(context.Background(), opts)
}

func (c *Client) GetFeedEntriesWithoutContext(feedID int, opts *EntryOptions) ([]Entry, *PaginationInfo, error) {
	return c.GetFeedEntries(context.Background(), feedID, opts)
}

func (c *Client) GetEntryWithoutContext(id int, opts *EntryOptions) (*Entry, error) {
	return c.GetEntry(context.Background(), id, opts)
}

func (c *Client) GetUnreadEntriesWithoutContext() ([]int, error) {
	return c.GetUnreadEntries(context.Background())
}

func (c *Client) MarkAsUnreadWithoutContext(entryIDs []int) ([]int, error) {
	return c.MarkAsUnread(context.Background(), entryIDs)
}

func (c *Client) MarkAsReadWithoutContext(entryIDs []int) ([]int, error) {
	return c.MarkAsRead(context.Background(), entryIDs)
}

func (c *Client) GetStarredEntriesWithoutContext() ([]int, error) {
	return c.GetStarredEntries(context.Background())
}

func (c *Client) StarEntriesWithoutContext(entryIDs []int) ([]int, error) {
	return c.StarEntries(context.Background(), entryIDs)
}

func (c *Client) UnstarEntriesWithoutContext(entryIDs []int) ([]int, error) {
	return c.UnstarEntries(context.Background(), entryIDs)
}

func (c *Client) GetEntriesByIDsWithoutContext(entryIDs []int, opts *EntryOptions) ([]Entry, error) {
	return c.GetEntriesByIDs(context.Background(), entryIDs, opts)
}
