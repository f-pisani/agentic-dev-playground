package feedbin

import (
	"context"
	"fmt"
)

// GetTaggings retrieves all taggings for the authenticated user.
// Returns a slice of taggings.
func (c *Client) GetTaggings(ctx context.Context) ([]Tagging, error) {
	var taggings []Tagging
	_, err := c.get(ctx, "taggings.json", nil, &taggings)
	if err != nil {
		return nil, err
	}

	return taggings, nil
}

// GetTagging retrieves a specific tagging by ID.
func (c *Client) GetTagging(ctx context.Context, id int) (*Tagging, error) {
	path := fmt.Sprintf("taggings/%d.json", id)

	var tagging Tagging
	_, err := c.get(ctx, path, nil, &tagging)
	if err != nil {
		return nil, err
	}

	return &tagging, nil
}

// CreateTagging creates a new tagging for a feed.
// Returns the created tagging or an error.
// If the tagging already exists, returns the existing tagging.
func (c *Client) CreateTagging(ctx context.Context, feedID int, name string) (*Tagging, error) {
	request := CreateTaggingRequest{
		FeedID: feedID,
		Name:   name,
	}

	var tagging Tagging
	resp, err := c.post(ctx, "taggings.json", nil, request, &tagging)
	if err != nil {
		return nil, err
	}

	// Handle 302 Found (existing tagging)
	if resp.StatusCode == 302 {
		// The tagging already exists, return it
		return &tagging, nil
	}

	return &tagging, nil
}

// DeleteTagging deletes a tagging by ID.
func (c *Client) DeleteTagging(ctx context.Context, id int) error {
	path := fmt.Sprintf("taggings/%d.json", id)
	_, err := c.delete(ctx, path, nil, nil)
	return err
}

// GetTaggingsByFeed retrieves all taggings for a specific feed.
func (c *Client) GetTaggingsByFeed(ctx context.Context, feedID int) ([]Tagging, error) {
	taggings, err := c.GetTaggings(ctx)
	if err != nil {
		return nil, err
	}

	var feedTaggings []Tagging
	for _, tagging := range taggings {
		if tagging.FeedID == feedID {
			feedTaggings = append(feedTaggings, tagging)
		}
	}

	return feedTaggings, nil
}

// GetTaggingsByName retrieves all taggings with a specific name.
func (c *Client) GetTaggingsByName(ctx context.Context, name string) ([]Tagging, error) {
	taggings, err := c.GetTaggings(ctx)
	if err != nil {
		return nil, err
	}

	var namedTaggings []Tagging
	for _, tagging := range taggings {
		if tagging.Name == name {
			namedTaggings = append(namedTaggings, tagging)
		}
	}

	return namedTaggings, nil
}

// DeleteTaggingsByFeed deletes all taggings for a specific feed.
func (c *Client) DeleteTaggingsByFeed(ctx context.Context, feedID int) error {
	taggings, err := c.GetTaggingsByFeed(ctx, feedID)
	if err != nil {
		return err
	}

	for _, tagging := range taggings {
		if err := c.DeleteTagging(ctx, tagging.ID); err != nil {
			return err
		}
	}

	return nil
}

// DeleteTaggingsByName deletes all taggings with a specific name.
func (c *Client) DeleteTaggingsByName(ctx context.Context, name string) error {
	taggings, err := c.GetTaggingsByName(ctx, name)
	if err != nil {
		return err
	}

	for _, tagging := range taggings {
		if err := c.DeleteTagging(ctx, tagging.ID); err != nil {
			return err
		}
	}

	return nil
}
