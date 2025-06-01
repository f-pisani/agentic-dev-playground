package feedbin

import (
	"net/http"
)

// ValidateAuth checks if the provided credentials are valid.
// It returns true if the credentials are valid, false otherwise.
func (c *Client) ValidateAuth() (bool, error) {
	resp, err := c.doRequest(http.MethodGet, "/authentication.json", nil, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Status code 200 means authentication was successful
	// Status code 401 means authentication failed
	return resp.StatusCode == http.StatusOK, nil
}
