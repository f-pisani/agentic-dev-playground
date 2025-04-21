package feedbin

import (
	"net/http"
)

// TagsService handles tag-related operations
type TagsService struct {
	client *Client
}

// List returns all tags
func (s *TagsService) List() ([]Tag, error) {
	req, err := s.client.newRequest(http.MethodGet, "/tags.json", nil)
	if err != nil {
		return nil, err
	}

	var tags []Tag
	_, err = s.client.do(req, &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
