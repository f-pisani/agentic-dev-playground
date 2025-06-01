package feedbin

import (
	"fmt"
	"net/http"
)

// GetTags retrieves all tags for the authenticated user.
func (c *Client) GetTags() ([]Tag, error) {
	// Make the request
	resp, err := c.doRequest(http.MethodGet, "/tags.json", nil, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var tags []Tag
	if err := parseResponse(resp, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}

// RenameTag renames a tag.
func (c *Client) RenameTag(oldName, newName string) error {
	// Create the request body
	reqBody := map[string]string{
		"old_name": oldName,
		"new_name": newName,
	}

	// Make the request
	resp, err := c.doRequest(http.MethodPost, "/tags.json", reqBody, nil)
	if err != nil {
		return err
	}

	// Check for success (200 OK)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to rename tag: %s", resp.Status)
	}

	return nil
}

// DeleteTag deletes a tag.
func (c *Client) DeleteTag(name string) error {
	// Create the request body
	reqBody := map[string]string{
		"name": name,
	}

	// Make the request
	resp, err := c.doRequest(http.MethodDelete, "/tags.json", reqBody, nil)
	if err != nil {
		return err
	}

	// Check for success (200 OK)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete tag: %s", resp.Status)
	}

	return nil
}
