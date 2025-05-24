package feedbin

import (
	"context"
)

// GetIcons retrieves feed icons for all subscribed feeds
func (c *Client) GetIcons(ctx context.Context) ([]Icon, error) {
	var icons []Icon
	_, err := c.makeRequest(ctx, "GET", "/icons.json", nil, &icons)
	if err != nil {
		return nil, err
	}

	return icons, nil
}

// GetIconsWithoutContext retrieves icons without context (uses background context)
func (c *Client) GetIconsWithoutContext() ([]Icon, error) {
	return c.GetIcons(context.Background())
}
