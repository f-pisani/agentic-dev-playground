package feedbin

import (
	"net/http"
)

// TagsService handles communication with the tags related
// endpoints of the Feedbin API
type TagsService struct {
	client *Client
}

// TagRenameOptions specifies the parameters to the
// TagsService.Rename method
type TagRenameOptions struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

// TagDeleteOptions specifies the parameters to the
// TagsService.Delete method
type TagDeleteOptions struct {
	Name string `json:"name"`
}

// Rename renames a tag
func (s *TagsService) Rename(opts *TagRenameOptions) ([]*Tagging, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "tags.json", opts)
	if err != nil {
		return nil, nil, err
	}

	var taggings []*Tagging
	resp, err := s.client.Do(req, &taggings)
	if err != nil {
		return nil, resp, err
	}

	return taggings, resp, nil
}

// Delete deletes a tag
func (s *TagsService) Delete(opts *TagDeleteOptions) ([]*Tagging, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodDelete, "tags.json", opts)
	if err != nil {
		return nil, nil, err
	}

	var taggings []*Tagging
	resp, err := s.client.Do(req, &taggings)
	if err != nil {
		return nil, resp, err
	}

	return taggings, resp, nil
}
