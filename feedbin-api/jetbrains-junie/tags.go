// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"fmt"
	"net/http"
)

// TagService handles communication with the tag related
// methods of the Feedbin API.
type TagService struct {
	client *Client
}

// GetTags retrieves all tags.
func (s *TagService) GetTags() ([]Tag, error) {
	req, err := s.client.NewRequest("GET", "/tags.json", nil)
	if err != nil {
		return nil, err
	}

	var tags []Tag
	_, err = s.client.Do(req, &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// RenameTag renames a tag.
func (s *TagService) RenameTag(oldName, newName string) error {
	body := map[string]interface{}{
		"old_name": oldName,
		"new_name": newName,
	}

	req, err := s.client.NewRequest("POST", "/tags.json", body)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// DeleteTag deletes a tag.
func (s *TagService) DeleteTag(name string) error {
	body := map[string]interface{}{
		"tag": name,
	}

	req, err := s.client.NewRequest("DELETE", "/tags.json", body)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}