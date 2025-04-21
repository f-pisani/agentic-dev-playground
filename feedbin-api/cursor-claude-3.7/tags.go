package feedbin

import (
	"net/http"
)

// TagsService handles communication with the tags related
// endpoints of the Feedbin API
type TagsService struct {
	client *Client
}

// List returns all tags
// https://github.com/feedbin/feedbin-api/blob/master/content/tags.md#get-tags
func (s *TagsService) List() ([]*Tag, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "tags.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var tags []*Tag
	resp, err := s.client.Do(req, &tags)
	if err != nil {
		return nil, resp, err
	}

	return tags, resp, nil
}
