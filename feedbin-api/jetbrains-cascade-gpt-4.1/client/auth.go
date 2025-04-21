package client

import (
	"net/http"
)

// CheckAuth verifies the credentials by calling /v2/authentication.json
func (c *Client) CheckAuth() error {
	req, err := c.newRequest(http.MethodGet, "/authentication.json", nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}
