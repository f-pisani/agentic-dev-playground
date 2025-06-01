package feedbin

import (
	"fmt"
	"net/http"
)

// GetTaggings retrieves all taggings for the authenticated user.
func (c *Client) GetTaggings() ([]Tagging, error) {
	// Make the request
	resp, err := c.doRequest(http.MethodGet, "/taggings.json", nil, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var taggings []Tagging
	if err := parseResponse(resp, &taggings); err != nil {
		return nil, err
	}

	return taggings, nil
}

// CreateTagging creates a new tagging (associates a tag with a feed).
func (c *Client) CreateTagging(feedID, tagID int) (*Tagging, error) {
	// Create the request body
	reqBody := TaggingRequest{
		FeedID: feedID,
		TagID:  tagID,
	}

	// Make the request
	resp, err := c.doRequest(http.MethodPost, "/taggings.json", reqBody, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var tagging Tagging
	if err := parseResponse(resp, &tagging); err != nil {
		return nil, err
	}

	return &tagging, nil
}

// DeleteTagging deletes a tagging by ID.
func (c *Client) DeleteTagging(id int) error {
	// Make the request
	resp, err := c.doRequest(http.MethodDelete, fmt.Sprintf("/taggings/%d.json", id), nil, nil)
	if err != nil {
		return err
	}

	// Check for success (204 No Content)
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
