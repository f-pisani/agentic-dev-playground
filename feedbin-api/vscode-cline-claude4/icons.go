package feedbin

import (
	"context"
	"fmt"
)

// GetIcons retrieves all icons for the authenticated user.
// Returns a slice of icons.
func (c *Client) GetIcons(ctx context.Context) ([]Icon, error) {
	var icons []Icon
	_, err := c.get(ctx, "icons.json", nil, &icons)
	if err != nil {
		return nil, err
	}

	return icons, nil
}

// GetIcon retrieves a specific icon by ID.
func (c *Client) GetIcon(ctx context.Context, id int) (*Icon, error) {
	path := fmt.Sprintf("icons/%d.json", id)

	var icon Icon
	_, err := c.get(ctx, path, nil, &icon)
	if err != nil {
		return nil, err
	}

	return &icon, nil
}

// GetIconByFeed retrieves the icon for a specific feed.
// Returns the first icon found for the feed, or nil if not found.
func (c *Client) GetIconByFeed(ctx context.Context, feedID int) (*Icon, error) {
	icons, err := c.GetIcons(ctx)
	if err != nil {
		return nil, err
	}

	for _, icon := range icons {
		if icon.FeedID == feedID {
			return &icon, nil
		}
	}

	return nil, fmt.Errorf("icon for feed %d not found", feedID)
}

// GetIconsByHost retrieves all icons for a specific host.
func (c *Client) GetIconsByHost(ctx context.Context, host string) ([]Icon, error) {
	icons, err := c.GetIcons(ctx)
	if err != nil {
		return nil, err
	}

	var hostIcons []Icon
	for _, icon := range icons {
		if icon.Host == host {
			hostIcons = append(hostIcons, icon)
		}
	}

	return hostIcons, nil
}
