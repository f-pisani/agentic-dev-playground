package feedbin

import (
	"context"
)

// ValidateCredentials validates the user's credentials by making a request to the authentication endpoint.
// Returns nil if credentials are valid, or an error if they are invalid or the request fails.
func (c *Client) ValidateCredentials(ctx context.Context) error {
	_, err := c.get(ctx, "authentication.json", nil, nil)
	return err
}
