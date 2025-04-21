package feedbin

import (
	"net/http"
)

// TagsService handles communication with the tags related
// methods of the Feedbin API.
type TagsService struct {
	client *Client
}

// RenameTagOptions specifies the parameters to the
// TagsService.Rename method.
type RenameTagOptions struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

// Rename renames a tag.
func (s *TagsService) Rename(opts *RenameTagOptions) ([]*Tagging, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/v2/tags.json", opts)
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

// DeleteTagOptions specifies the parameters to the
// TagsService.Delete method.
type DeleteTagOptions struct {
	Name string `json:"name"`
}

// Delete deletes a tag.
func (s *TagsService) Delete(opts *DeleteTagOptions) ([]*Tagging, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodDelete, "/v2/tags.json", opts)
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
