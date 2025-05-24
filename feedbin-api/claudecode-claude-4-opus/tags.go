package feedbin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetTags retrieves all tags
func (c *Client) GetTags() ([]Tag, error) {
	var tags []Tag
	err := c.get("/tags.json", nil, &tags)
	return tags, err
}

// RenameTag renames an existing tag
func (c *Client) RenameTag(oldName, newName string) ([]Tagging, error) {
	body := struct {
		OldName string `json:"old_name"`
		NewName string `json:"new_name"`
	}{
		OldName: oldName,
		NewName: newName,
	}

	var taggings []Tagging
	err := c.post("/tags.json", body, &taggings)
	if err != nil {
		return nil, err
	}
	return taggings, nil
}

// DeleteTag deletes a tag by name
func (c *Client) DeleteTag(name string) ([]Tagging, error) {
	body := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}

	var taggings []Tagging
	resp, err := c.request(http.MethodDelete, "/tags.json", body, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&taggings); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return taggings, nil
}

// GetTagByName retrieves a tag by its name
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

	return nil, &NotFoundError{
		APIError: &APIError{
			StatusCode: 404,
			Status:     "Not Found",
			Message:    fmt.Sprintf("tag not found: %s", name),
		},
	}
}
