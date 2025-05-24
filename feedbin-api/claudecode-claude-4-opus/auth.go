package feedbin

// CheckAuthentication verifies the provided credentials
// Returns nil if authentication is successful, or an error if it fails
func (c *Client) CheckAuthentication() error {
	var auth Authentication
	err := c.get("/authentication.json", nil, &auth)
	if err != nil {
		return err
	}
	return nil
}

// GetAuthenticatedUser returns the email of the authenticated user
func (c *Client) GetAuthenticatedUser() (string, error) {
	var auth Authentication
	err := c.get("/authentication.json", nil, &auth)
	if err != nil {
		return "", err
	}
	return c.Email, nil
}
