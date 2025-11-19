package feedbin

import "context"

// GetUpdatedEntries fetches the IDs of updated entries.
func (c *Client) GetUpdatedEntries(ctx context.Context, opt *PaginationOptions) ([]int64, *Response, error) {
	u, err := addOptions("updated_entries.json", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.NewRequest("GET", u, nil)
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

type updatedEntriesRequest struct {
	UpdatedEntries []int64 `json:"updated_entries"`
}

// DeleteUpdatedEntries marks updated entries as read.
func (c *Client) DeleteUpdatedEntries(ctx context.Context, ids []int64) ([]int64, *Response, error) {
	reqBody := &updatedEntriesRequest{UpdatedEntries: ids}
	req, err := c.NewRequest("DELETE", "updated_entries.json", reqBody)
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
