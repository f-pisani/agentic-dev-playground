package client

import (
	"net/http"

	"feedbin-api/jetbrains-cascade-gpt-4.1/models"
)

// GetTags fetches all tags for the user.
func (c *Client) GetTags() ([]models.Tag, error) {
	req, err := c.newRequest(http.MethodGet, "/tags.json", nil)
	if err != nil {
		return nil, err
	}
	var tags models.TagsResponse
	err = c.do(req, &tags)
	if err != nil {
		return nil, err
	}
	return tags, nil
}
