package feedbin

import (
	"net/http"
)

// TagsService handles operations related to tags
type TagsService struct {
	client *Client
}

// Rename renames a tag
// Returns the updated taggings after the rename
func (s *TagsService) Rename(oldName, newName string) ([]Tagging, error) {
	params := &TagParams{
		OldName: oldName,
		NewName: newName,
	}

	req, err := s.client.NewRequest(http.MethodPost, "/tags.json", params)
	if err != nil {
		return nil, err
	}

	var taggings []Tagging
	_, err = s.client.Do(req, &taggings)
	if err != nil {
		return nil, err
	}

	return taggings, nil
}

// Delete deletes a tag
// Returns the updated taggings after the delete
func (s *TagsService) Delete(name string) ([]Tagging, error) {
	params := &TagDeleteParams{
		Name: name,
	}

	req, err := s.client.NewRequest(http.MethodDelete, "/tags.json", params)
	if err != nil {
		return nil, err
	}

	var taggings []Tagging
	_, err = s.client.Do(req, &taggings)
	if err != nil {
		return nil, err
	}

	return taggings, nil
}
