package feedbin

import (
	"context"
	"fmt"
)

// GetTaggings fetches all taggings.
func (c *Client) GetTaggings(ctx context.Context) ([]Tagging, *Response, error) {
	req, err := c.NewRequest("GET", "taggings.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var taggings []Tagging
	resp, err := c.Do(ctx, req, &taggings)
	if err != nil {
		return nil, resp, err
	}

	return taggings, resp, nil
}

// GetTagging fetches a single tagging by ID.
func (c *Client) GetTagging(ctx context.Context, id int64) (*Tagging, *Response, error) {
	u := fmt.Sprintf("taggings/%d.json", id)
	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var tagging Tagging
	resp, err := c.Do(ctx, req, &tagging)
	if err != nil {
		return nil, resp, err
	}

	return &tagging, resp, nil
}

type taggingCreateRequest struct {
	FeedID int64  `json:"feed_id"`
	Name   string `json:"name"`
}

// CreateTagging creates a new tagging.
func (c *Client) CreateTagging(ctx context.Context, feedID int64, name string) (*Tagging, *Response, error) {
	reqBody := &taggingCreateRequest{FeedID: feedID, Name: name}
	req, err := c.NewRequest("POST", "taggings.json", reqBody)
	if err != nil {
		return nil, nil, err
	}

	var tagging Tagging
	resp, err := c.Do(ctx, req, &tagging)
	if err != nil {
		return nil, resp, err
	}

	return &tagging, resp, nil
}

// DeleteTagging deletes a tagging.
func (c *Client) DeleteTagging(ctx context.Context, id int64) (*Response, error) {
	u := fmt.Sprintf("taggings/%d.json", id)
	req, err := c.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(ctx, req, nil)
}
