package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"feedbin-api/jetbrains-cascade-gpt-4.1/models"
)

// PageRequest is the payload for creating a new page.
type PageRequest struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

// CreatePage creates a new page (entry) from a URL.
func (c *Client) CreatePage(reqBody PageRequest) (*models.Entry, error) {
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest(http.MethodPost, "/pages.json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var entry models.Entry
	err = c.do(req, &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}
