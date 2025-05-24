package feedbin

// GetIcons retrieves all available favicon data
func (c *Client) GetIcons() ([]Icon, error) {
	var icons []Icon
	err := c.get("/icons.json", nil, &icons)
	return icons, err
}

// GetIconByHost retrieves favicon data for a specific host
func (c *Client) GetIconByHost(host string) (*Icon, error) {
	icons, err := c.GetIcons()
	if err != nil {
		return nil, err
	}

	for _, icon := range icons {
		if icon.Host == host {
			return &icon, nil
		}
	}

	return nil, &NotFoundError{
		APIError: &APIError{
			StatusCode: 404,
			Status:     "Not Found",
			Message:    "icon not found for host: " + host,
		},
	}
}
