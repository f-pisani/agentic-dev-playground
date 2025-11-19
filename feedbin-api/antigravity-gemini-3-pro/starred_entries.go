package feedbin

import "context"

// GetStarredEntries fetches the IDs of all starred entries.
func (c *Client) GetStarredEntries(ctx context.Context) ([]int64, *Response, error) {
	req, err := c.NewRequest("GET", "starred_entries.json", nil)
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

type starredEntriesRequest struct {
	StarredEntries []int64 `json:"starred_entries"`
}

// CreateStarredEntries stars entries.
func (c *Client) CreateStarredEntries(ctx context.Context, ids []int64) ([]int64, *Response, error) {
	reqBody := &starredEntriesRequest{StarredEntries: ids}
	req, err := c.NewRequest("POST", "starred_entries.json", reqBody)
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

// DeleteStarredEntries unstars entries.
func (c *Client) DeleteStarredEntries(ctx context.Context, ids []int64) ([]int64, *Response, error) {
	reqBody := &starredEntriesRequest{StarredEntries: ids}
	req, err := c.NewRequest("DELETE", "starred_entries.json", reqBody)
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
