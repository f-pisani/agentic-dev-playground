package client

import (
	"net/http"

	"feedbin-api/jetbrains-cascade-gpt-4.1/models"
)

// GetIcons fetches all feed icons for the user.
func (c *Client) GetIcons() ([]models.Icon, error) {
	req, err := c.newRequest(http.MethodGet, "/icons.json", nil)
	if err != nil {
		return nil, err
	}
	var icons models.IconsResponse
	err = c.do(req, &icons)
	if err != nil {
		return nil, err
	}
	return icons, nil
}
