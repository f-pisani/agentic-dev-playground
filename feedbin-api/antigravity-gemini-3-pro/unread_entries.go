package feedbin

import "context"

// GetUnreadEntries fetches the IDs of all unread entries.
func (c *Client) GetUnreadEntries(ctx context.Context) ([]int64, *Response, error) {
	req, err := c.NewRequest("GET", "unread_entries.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var ids []int64
	resp, err := c.Do(ctx, req, &ids)
	if err != nil {
		return nil, resp, err
	}

	return ids, resp, nil
}

type unreadEntriesRequest struct {
	UnreadEntries []int64 `json:"unread_entries"`
}

// CreateUnreadEntries marks entries as unread.
func (c *Client) CreateUnreadEntries(ctx context.Context, ids []int64) ([]int64, *Response, error) {
	reqBody := &unreadEntriesRequest{UnreadEntries: ids}
	req, err := c.NewRequest("POST", "unread_entries.json", reqBody)
	if err != nil {
		return nil, nil, err
	}

	var createdIDs []int64
	resp, err := c.Do(ctx, req, &createdIDs)
	if err != nil {
		return nil, resp, err
	}

	return createdIDs, resp, nil
}

// DeleteUnreadEntries marks entries as read.
func (c *Client) DeleteUnreadEntries(ctx context.Context, ids []int64) ([]int64, *Response, error) {
	reqBody := &unreadEntriesRequest{UnreadEntries: ids}
	req, err := c.NewRequest("DELETE", "unread_entries.json", reqBody)
	if err != nil {
		return nil, nil, err
	}

	var deletedIDs []int64
	resp, err := c.Do(ctx, req, &deletedIDs)
	if err != nil {
		return nil, resp, err
	}

	return deletedIDs, resp, nil
}
