package feedbin

import (
	"context"
)

// RenameTag renames a tag from oldName to newName.
// Returns the updated taggings after the rename operation.
func (c *Client) RenameTag(ctx context.Context, oldName, newName string) ([]Tagging, error) {
	request := RenameTagRequest{
		OldName: oldName,
		NewName: newName,
	}

	var taggings []Tagging
	_, err := c.post(ctx, "tags.json", nil, request, &taggings)
	if err != nil {
		return nil, err
	}

	return taggings, nil
}

// DeleteTag deletes a tag by name.
// Returns the remaining taggings after the delete operation.
func (c *Client) DeleteTag(ctx context.Context, name string) ([]Tagging, error) {
	request := DeleteTagRequest{
		Name: name,
	}

	var taggings []Tagging
	_, err := c.delete(ctx, "tags.json", nil, request)
	if err != nil {
		return nil, err
	}

	// Since DELETE doesn't return a body, we need to get the updated taggings
	taggings, err = c.GetTaggings(ctx)
	if err != nil {
		return nil, err
	}

	return taggings, nil
}

// GetAllTagNames retrieves all unique tag names for the authenticated user.
// This is a convenience method that extracts unique tag names from all taggings.
func (c *Client) GetAllTagNames(ctx context.Context) ([]string, error) {
	taggings, err := c.GetTaggings(ctx)
	if err != nil {
		return nil, err
	}

	// Use a map to collect unique tag names
	tagMap := make(map[string]bool)
	for _, tagging := range taggings {
		tagMap[tagging.Name] = true
	}

	// Convert map keys to slice
	var tagNames []string
	for name := range tagMap {
		tagNames = append(tagNames, name)
	}

	return tagNames, nil
}

// GetFeedsByTag retrieves all feed IDs that have a specific tag.
// This is a convenience method that filters taggings by tag name.
func (c *Client) GetFeedsByTag(ctx context.Context, tagName string) ([]int, error) {
	taggings, err := c.GetTaggingsByName(ctx, tagName)
	if err != nil {
		return nil, err
	}

	var feedIDs []int
	for _, tagging := range taggings {
		feedIDs = append(feedIDs, tagging.FeedID)
	}

	return feedIDs, nil
}

// GetTagsByFeed retrieves all tag names for a specific feed.
// This is a convenience method that extracts tag names from feed taggings.
func (c *Client) GetTagsByFeed(ctx context.Context, feedID int) ([]string, error) {
	taggings, err := c.GetTaggingsByFeed(ctx, feedID)
	if err != nil {
		return nil, err
	}

	var tagNames []string
	for _, tagging := range taggings {
		tagNames = append(tagNames, tagging.Name)
	}

	return tagNames, nil
}
