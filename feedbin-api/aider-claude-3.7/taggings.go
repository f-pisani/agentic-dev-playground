package feedbin

import (
	"fmt"
	"net/http"
)

// GetTaggings retrieves all taggings
func (c *Client) GetTaggings() ([]Tagging, error) {
	req, err := c.NewRequest(http.MethodGet, "/taggings.json", nil)
	if err != nil {
		return nil, err
	}
	
	var taggings []Tagging
	_, err = c.Do(req, &taggings)
	if err != nil {
		return nil, err
	}
	
	return taggings, nil
}

// CreateTagging creates a new tagging (adds a tag to a feed)
func (c *Client) CreateTagging(feedID int64, name string) (*Tagging, error) {
	taggingReq := &TaggingRequest{
		FeedID: feedID,
		Name:   name,
	}
	
	req, err := c.NewRequest(http.MethodPost, "/taggings.json", taggingReq)
	if err != nil {
		return nil, err
	}
	
	tagging := new(Tagging)
	_, err = c.Do(req, tagging)
	if err != nil {
		return nil, err
	}
	
	return tagging, nil
}

// DeleteTagging deletes a tagging
func (c *Client) DeleteTagging(id int64) error {
	path := fmt.Sprintf("/taggings/%d.json", id)
	req, err := c.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	
	_, err = c.Do(req, nil)
	return err
}

// GetTaggingsByFeed retrieves all taggings for a specific feed
func (c *Client) GetTaggingsByFeed(feedID int64) ([]Tagging, error) {
	taggings, err := c.GetTaggings()
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

// GetTaggingsByTag retrieves all taggings for a specific tag
func (c *Client) GetTaggingsByTag(tagID int64) ([]Tagging, error) {
	taggings, err := c.GetTaggings()
	if err != nil {
		return nil, err
	}
	
	var tagTaggings []Tagging
	for _, tagging := range taggings {
		if tagging.TagID == tagID {
			tagTaggings = append(tagTaggings, tagging)
		}
	}
	
	return tagTaggings, nil
}
