package feedbin

import (
	"fmt"
	"net/http"
)

// GetIcon retrieves the icon for a feed.
func (c *Client) GetIcon(feedID int) (*Icon, error) {
	// Make the request
	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("/feeds/%d/icon.json", feedID), nil, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var icon Icon
	if err := parseResponse(resp, &icon); err != nil {
		return nil, err
	}

	return &icon, nil
}
