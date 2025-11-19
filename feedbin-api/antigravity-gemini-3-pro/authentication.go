package feedbin

import "context"

// CheckAuthentication checks if the provided credentials are valid.
// It returns nil if authentication is successful, or an error otherwise.
func (c *Client) CheckAuthentication(ctx context.Context) error {
	req, err := c.NewRequest("GET", "authentication.json", nil)
	if err != nil {
		return err
	}

	_, err = c.Do(ctx, req, nil)
	return err
}
