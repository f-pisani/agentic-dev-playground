package feedbinapi

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// TagsService handles operations related to tags.
type TagsService struct {
	client *Client
}

// NewTagsService creates a new service for tag related operations.
func NewTagsService(client *Client) *TagsService {
	return &TagsService{client: client}
}

// TagListOptions specifies optional parameters for listing tags.
// The API doc for "GET /v2/tags.json" does not specify query parameters.
type TagListOptions struct {
	// No options defined in spec, but can be added for consistency.
	// For example, if sorting or pagination were ever added.
}

// List retrieves all unique tags.
// Each object is a simple tag with an ID and name.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/tags.md#get-tags
func (s *TagsService) List(opts *TagListOptions) ([]Tag, *http.Response, error) {
	path := "tags.json"
	if opts != nil {
		v, err := query.Values(opts)
		if err != nil {
			return nil, nil, err
		}
		if params := v.Encode(); params != "" {
			path = fmt.Sprintf("%s?%s", path, params)
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var tags []Tag
	resp, err := s.client.Do(req, &tags)
	if err != nil {
		return nil, resp, err
	}
	return tags, resp, nil
}

// Delete removes a tag and all its associated taggings.
// The API specifies the tag ID in the path: DELETE /v2/tags/{tag_id}.json
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/tags.md#delete-tag
func (s *TagsService) Delete(tagID int64) (*http.Response, error) {
	if tagID == 0 {
		return nil, fmt.Errorf("tagID is required to delete a tag")
	}
	path := fmt.Sprintf("tags/%d.json", tagID)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	// Successful deletion should return 204 No Content.
	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
