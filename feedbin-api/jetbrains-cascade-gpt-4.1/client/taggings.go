package client

import (
	"net/http"

	"feedbin-api/jetbrains-cascade-gpt-4.1/models"
)

// GetTaggings fetches all taggings for the user.
func (c *Client) GetTaggings() ([]models.Tagging, error) {
	req, err := c.newRequest(http.MethodGet, "/taggings.json", nil)
	if err != nil {
		return nil, err
	}
	var taggings models.TaggingsResponse
	err = c.do(req, &taggings)
	if err != nil {
		return nil, err
	}
	return taggings, nil
}
