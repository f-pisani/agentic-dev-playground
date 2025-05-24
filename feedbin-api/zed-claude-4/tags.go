package feedbin

import (
	"context"
	"fmt"
)

// GetTaggings retrieves all taggings for the authenticated user
func (c *Client) GetTaggings(ctx context.Context) ([]Tagging, error) {
	var taggings []Tagging
	_, err := c.makeRequest(ctx, "GET", "/taggings.json", nil, &taggings)
	if err != nil {
		return nil, err
	}

	return taggings, nil
}

// GetTagging retrieves a specific tagging by ID
func (c *Client) GetTagging(ctx context.Context, id int) (*Tagging, error) {
	path := fmt.Sprintf("/taggings/%d.json", id)

	var tagging Tagging
	_, err := c.makeRequest(ctx, "GET", path, nil, &tagging)
	if err != nil {
		return nil, err
	}

	return &tagging, nil
}

// CreateTagging creates a new tagging for a feed
func (c *Client) CreateTagging(ctx context.Context, feedID int, name string) (*Tagging, error) {
	if feedID <= 0 {
		return nil, &ValidationError{
			Field:   "feed_id",
			Message: "feed ID must be greater than 0",
		}
	}

	if name == "" {
		return nil, &ValidationError{
			Field:   "name",
			Message: "tag name is required",
		}
	}

	req := &CreateTaggingRequest{
		FeedID: feedID,
		Name:   name,
	}

	var tagging Tagging
	_, err := c.makeRequest(ctx, "POST", "/taggings.json", req, &tagging)
	if err != nil {
		return nil, err
	}

	return &tagging, nil
}

// DeleteTagging deletes a tagging by ID
func (c *Client) DeleteTagging(ctx context.Context, id int) error {
	path := fmt.Sprintf("/taggings/%d.json", id)

	_, err := c.makeRequest(ctx, "DELETE", path, nil, nil)
	return err
}

// RenameTag renames a tag across all taggings
func (c *Client) RenameTag(ctx context.Context, oldName, newName string) ([]Tagging, error) {
	if oldName == "" {
		return nil, &ValidationError{
			Field:   "old_name",
			Message: "old tag name is required",
		}
	}

	if newName == "" {
		return nil, &ValidationError{
			Field:   "new_name",
			Message: "new tag name is required",
		}
	}

	req := &RenameTagRequest{
		OldName: oldName,
		NewName: newName,
	}

	var taggings []Tagging
	_, err := c.makeRequest(ctx, "POST", "/tags.json", req, &taggings)
	if err != nil {
		return nil, err
	}

	return taggings, nil
}

// DeleteTag deletes a tag across all taggings
func (c *Client) DeleteTag(ctx context.Context, name string) ([]Tagging, error) {
	if name == "" {
		return nil, &ValidationError{
			Field:   "name",
			Message: "tag name is required",
		}
	}

	req := &DeleteTagRequest{
		Name: name,
	}

	var taggings []Tagging
	_, err := c.makeRequest(ctx, "DELETE", "/tags.json", req, &taggings)
	if err != nil {
		return nil, err
	}

	return taggings, nil
}

// GetTaggingsByFeed retrieves all taggings for a specific feed
func (c *Client) GetTaggingsByFeed(ctx context.Context, feedID int) ([]Tagging, error) {
	taggings, err := c.GetTaggings(ctx)
	if err != nil {
		return nil, err
	}

	var result []Tagging
	for _, tagging := range taggings {
		if tagging.FeedID == feedID {
			result = append(result, tagging)
		}
	}

	return result, nil
}

// GetTaggingsByName retrieves all taggings with a specific tag name
func (c *Client) GetTaggingsByName(ctx context.Context, name string) ([]Tagging, error) {
	taggings, err := c.GetTaggings(ctx)
	if err != nil {
		return nil, err
	}

	var result []Tagging
	for _, tagging := range taggings {
		if tagging.Name == name {
			result = append(result, tagging)
		}
	}

	return result, nil
}

// GetUniqueTagNames retrieves all unique tag names
func (c *Client) GetUniqueTagNames(ctx context.Context) ([]string, error) {
	taggings, err := c.GetTaggings(ctx)
	if err != nil {
		return nil, err
	}

	nameSet := make(map[string]bool)
	for _, tagging := range taggings {
		nameSet[tagging.Name] = true
	}

	var names []string
	for name := range nameSet {
		names = append(names, name)
	}

	return names, nil
}

// Convenience methods without context (use background context)

func (c *Client) GetTaggingsWithoutContext() ([]Tagging, error) {
	return c.GetTaggings(context.Background())
}

func (c *Client) GetTaggingWithoutContext(id int) (*Tagging, error) {
	return c.GetTagging(context.Background(), id)
}

func (c *Client) CreateTaggingWithoutContext(feedID int, name string) (*Tagging, error) {
	return c.CreateTagging(context.Background(), feedID, name)
}

func (c *Client) DeleteTaggingWithoutContext(id int) error {
	return c.DeleteTagging(context.Background(), id)
}

func (c *Client) RenameTagWithoutContext(oldName, newName string) ([]Tagging, error) {
	return c.RenameTag(context.Background(), oldName, newName)
}

func (c *Client) DeleteTagWithoutContext(name string) ([]Tagging, error) {
	return c.DeleteTag(context.Background(), name)
}

func (c *Client) GetTaggingsByFeedWithoutContext(feedID int) ([]Tagging, error) {
	return c.GetTaggingsByFeed(context.Background(), feedID)
}

func (c *Client) GetTaggingsByNameWithoutContext(name string) ([]Tagging, error) {
	return c.GetTaggingsByName(context.Background(), name)
}

func (c *Client) GetUniqueTagNamesWithoutContext() ([]string, error) {
	return c.GetUniqueTagNames(context.Background())
}
