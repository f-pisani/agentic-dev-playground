package feedbin

import "context"

// GetIcons fetches icons for all subscribed feeds.
func (c *Client) GetIcons(ctx context.Context) ([]Icon, *Response, error) {
	req, err := c.NewRequest("GET", "icons.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var icons []Icon
	resp, err := c.Do(ctx, req, &icons)
	if err != nil {
		return nil, resp, err
	}

	return icons, resp, nil
}
