package feedbinapi

import (
	"context"
	"fmt"
	"net/http"
)

// TagsService handles operations related to renaming or deleting tags globally.
type TagsService struct {
	client *Client
}

// RenameTag changes the name of a tag across all feeds.
// The response is an array of taggings affected by the rename.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/tags.md#post-v2tagsjson
func (s *TagsService) RenameTag(ctx context.Context, oldName string, newName string) ([]Tagging, *Response, error) {
	if oldName == "" || newName == "" {
		return nil, nil, fmt.Errorf("oldName and newName cannot be empty")
	}
	if oldName == newName {
		return nil, nil, fmt.Errorf("oldName and newName cannot be the same")
	}

	path := "tags.json"
	requestBody := &RenameTagRequest{OldName: oldName, NewName: newName}

	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create RenameTag request: %w", err)
	}

	var affectedTaggings []Tagging
	resp, err := s.client.do(req, &affectedTaggings)
	if err != nil {
		return nil, resp, err
	}
	return affectedTaggings, resp, nil
}

// DeleteTag removes a tag from all feeds it was applied to.
// The response is an array of remaining taggings (potentially, or needs clarification from docs).
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/tags.md#delete-v2tagsjson
func (s *TagsService) DeleteTag(ctx context.Context, name string) ([]Tagging, *Response, error) {
	if name == "" {
		return nil, nil, fmt.Errorf("name of tag to delete cannot be empty")
	}

	path := "tags.json"
	requestBody := &DeleteTagRequest{Name: name}

	// Note: The API uses DELETE /v2/tags.json but with a request body.
	// http.NewRequestWithContext will allow a body with DELETE.
	req, err := s.client.newRequest(ctx, http.MethodDelete, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create DeleteTag request: %w", err)
	}

	var remainingTaggings []Tagging // Docs say "new array of taggings after the delete"
	resp, err := s.client.do(req, &remainingTaggings)
	if err != nil {
		return nil, resp, err
	}
	return remainingTaggings, resp, nil
}
