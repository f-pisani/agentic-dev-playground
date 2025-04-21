package feedbin

import (
	"fmt"
	"net/http"
)

// GetIcon retrieves a specific icon by ID
func (c *Client) GetIcon(id int64) (*Favicon, error) {
	path := fmt.Sprintf("/v2/icons/%d.json", id)
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	
	icon := new(Favicon)
	_, err = c.Do(req, icon)
	if err != nil {
		return nil, err
	}
	
	return icon, nil
}

// GetFeedIcon retrieves the icon for a specific feed
func (c *Client) GetFeedIcon(feedID int64) (*Favicon, error) {
	// First get the subscription to get the favicon ID
	subscription, err := c.GetSubscription(feedID)
	if err != nil {
		return nil, err
	}
	
	if subscription.Favicon == nil {
		return nil, fmt.Errorf("feed %d has no favicon", feedID)
	}
	
	// Then get the icon
	return c.GetIcon(subscription.Favicon.ID)
}
