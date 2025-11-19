package feedbin

import "context"

// GetRecentlyReadEntries fetches the IDs of recently read entries.
func (c *Client) GetRecentlyReadEntries(ctx context.Context) ([]int64, *Response, error) {
	req, err := c.NewRequest("GET", "recently_read_entries.json", nil)
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

type recentlyReadEntriesRequest struct {
	RecentlyReadEntries []int64 `json:"recently_read_entries"`
}

// CreateRecentlyReadEntries creates records for recently read entries.
func (c *Client) CreateRecentlyReadEntries(ctx context.Context, ids []int64) ([]int64, *Response, error) {
	reqBody := &recentlyReadEntriesRequest{RecentlyReadEntries: ids}
	req, err := c.NewRequest("POST", "recently_read_entries.json", reqBody)
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
