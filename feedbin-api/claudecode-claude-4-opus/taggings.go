package feedbin

import (
	"fmt"
)

// GetTaggings retrieves all taggings
func (c *Client) GetTaggings() ([]Tagging, error) {
	var taggings []Tagging
	err := c.get("/taggings.json", nil, &taggings)
	return taggings, err
}

// GetTagging retrieves a single tagging by ID
func (c *Client) GetTagging(id int) (*Tagging, error) {
	var tagging Tagging
	path := fmt.Sprintf("/taggings/%d.json", id)
	err := c.get(path, nil, &tagging)
	if err != nil {
		return nil, err
	}
	return &tagging, nil
}

// CreateTagging creates a new tagging (associates a feed with a tag)
func (c *Client) CreateTagging(feedID int, tagName string) (*Tagging, error) {
	body := struct {
		FeedID int    `json:"feed_id"`
		Name   string `json:"name"`
	}{
		FeedID: feedID,
		Name:   tagName,
	}

	var tagging Tagging
	err := c.post("/taggings.json", body, &tagging)
	if err != nil {
		return nil, err
	}
	return &tagging, nil
}

// DeleteTagging deletes a tagging
func (c *Client) DeleteTagging(id int) error {
	path := fmt.Sprintf("/taggings/%d.json", id)
	return c.delete(path, nil)
}

// GetTaggingsByFeed retrieves all taggings for a specific feed
func (c *Client) GetTaggingsByFeed(feedID int) ([]Tagging, error) {
	allTaggings, err := c.GetTaggings()
	if err != nil {
		return nil, err
	}

	var feedTaggings []Tagging
	for _, tagging := range allTaggings {
		if tagging.FeedID == feedID {
			feedTaggings = append(feedTaggings, tagging)
		}
	}

	return feedTaggings, nil
}

// GetTaggingsByName retrieves all taggings with a specific tag name
func (c *Client) GetTaggingsByName(tagName string) ([]Tagging, error) {
	allTaggings, err := c.GetTaggings()
	if err != nil {
		return nil, err
	}

	var namedTaggings []Tagging
	for _, tagging := range allTaggings {
		if tagging.Name == tagName {
			namedTaggings = append(namedTaggings, tagging)
		}
	}

	return namedTaggings, nil
}

// TagFeed applies a tag to a feed
func (c *Client) TagFeed(feedID int, tagName string) (*Tagging, error) {
	return c.CreateTagging(feedID, tagName)
}

// UntagFeed removes a tag from a feed
func (c *Client) UntagFeed(feedID int, tagName string) error {
	taggings, err := c.GetTaggingsByFeed(feedID)
	if err != nil {
		return err
	}

	for _, tagging := range taggings {
		if tagging.Name == tagName {
			return c.DeleteTagging(tagging.ID)
		}
	}

	return &NotFoundError{
		APIError: &APIError{
			StatusCode: 404,
			Status:     "Not Found",
			Message:    fmt.Sprintf("tagging not found for feed %d with tag %s", feedID, tagName),
		},
	}
}
