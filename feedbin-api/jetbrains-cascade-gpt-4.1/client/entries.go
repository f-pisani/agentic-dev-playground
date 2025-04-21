package client

import (
	"net/http"
	"net/url"
	"strconv"

	"feedbin-api/jetbrains-cascade-gpt-4.1/models"
)

// EntriesOptions holds query params for entries endpoint.
type EntriesOptions struct {
	Page    int
	PerPage int
}

// GetEntries fetches entries with pagination.
func (c *Client) GetEntries(opt EntriesOptions) ([]models.Entry, error) {
	v := url.Values{}
	if opt.Page > 0 {
		v.Set("page", strconv.Itoa(opt.Page))
	}
	if opt.PerPage > 0 {
		v.Set("per_page", strconv.Itoa(opt.PerPage))
	}
	path := "/entries.json"
	if len(v) > 0 {
		path += "?" + v.Encode()
	}
	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var entries models.EntriesResponse
	err = c.do(req, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
