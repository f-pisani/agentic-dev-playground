package feedbin

import (
	"fmt"
	"net/http"
)

// GetTags retrieves all tags
func (c *Client) GetTags() ([]Tag, error) {
	req, err := c.NewRequest(http.MethodGet, "/tags.json", nil)
	if err != nil {
		return nil, err
	}
	
	var tags []Tag
	_, err = c.Do(req, &tags)
	if err != nil {
		return nil, err
	}
	
	return tags, nil
}

// DeleteTag deletes a tag and all associated taggings
func (c *Client) DeleteTag(id int64) error {
	path := fmt.Sprintf("/tags/%d.json", id)
	req, err := c.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	
	_, err = c.Do(req, nil)
	return err
}

// GetTagByName finds a tag by its name
func (c *Client) GetTagByName(name string) (*Tag, error) {
	tags, err := c.GetTags()
	if err != nil {
		return nil, err
	}
	
	for _, tag := range tags {
		if tag.Name == name {
			return &tag, nil
		}
	}
	
	return nil, fmt.Errorf("tag not found: %s", name)
}

// GetFeedsByTag retrieves all feeds with a specific tag
func (c *Client) GetFeedsByTag(tagID int64) ([]int64, error) {
	taggings, err := c.GetTaggingsByTag(tagID)
	if err != nil {
		return nil, err
	}
	
	var feedIDs []int64
	for _, tagging := range taggings {
		feedIDs = append(feedIDs, tagging.FeedID)
	}
	
	return feedIDs, nil
}

// GetTagsWithCounts retrieves all tags with their feed counts
func (c *Client) GetTagsWithCounts() ([]Tag, error) {
	tags, err := c.GetTags()
	if err != nil {
		return nil, err
	}
	
	taggings, err := c.GetTaggings()
	if err != nil {
		return nil, err
	}
	
	// Count feeds per tag
	tagCounts := make(map[int64]int)
	for _, tagging := range taggings {
		tagCounts[tagging.TagID]++
	}
	
	// Add counts to tags
	for i := range tags {
		tags[i].Count = tagCounts[tags[i].ID]
	}
	
	return tags, nil
}
