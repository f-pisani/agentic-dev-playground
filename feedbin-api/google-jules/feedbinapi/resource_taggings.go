package feedbinapi

import (
	"context"
	"fmt"
	"net/http"
)

// TaggingsService handles operations related to feed taggings.
type TaggingsService struct {
	client *Client
}

// ListTaggings retrieves all taggings for the authenticated user.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/taggings.md#get-taggings
func (s *TaggingsService) ListTaggings(ctx context.Context) ([]Tagging, *Response, error) {
	path := "taggings.json"
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ListTaggings request: %w", err)
	}

	var taggings []Tagging
	resp, err := s.client.do(req, &taggings)
	if err != nil {
		return nil, resp, err
	}
	return taggings, resp, nil
}

// GetTagging retrieves a single tagging by its ID.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/taggings.md#get-tagging
func (s *TaggingsService) GetTagging(ctx context.Context, taggingID int64) (*Tagging, *Response, error) {
	if taggingID <= 0 {
		return nil, nil, fmt.Errorf("taggingID must be a positive integer")
	}
	path := fmt.Sprintf("taggings/%d.json", taggingID)
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GetTagging request: %w", err)
	}

	var tagging Tagging
	resp, err := s.client.do(req, &tagging)
	if err != nil {
		return nil, resp, err
	}
	return &tagging, resp, nil
}

// CreateTagging creates a new tagging, assigning a name (tag) to a specific feed.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/taggings.md#create-tagging
func (s *TaggingsService) CreateTagging(ctx context.Context, feedID int64, name string) (*Tagging, *Response, error) {
	if feedID <= 0 {
		return nil, nil, fmt.Errorf("feedID must be a positive integer")
	}
	if name == "" {
		return nil, nil, fmt.Errorf("name for tagging cannot be empty")
	}

	path := "taggings.json"
	requestBody := &CreateTaggingRequest{FeedID: feedID, Name: name}

	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CreateTagging request: %w", err)
	}

	var createdTagging Tagging
	// API returns 201 Created or 302 Found. `do` method handles unmarshalling for these.
	resp, err := s.client.do(req, &createdTagging)
	if err != nil {
		return nil, resp, err
	}
	return &createdTagging, resp, nil
}

// DeleteTagging removes a tagging by its ID.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/taggings.md#delete-tagging
func (s *TaggingsService) DeleteTagging(ctx context.Context, taggingID int64) (*Response, error) {
	if taggingID <= 0 {
		return nil, fmt.Errorf("taggingID must be a positive integer")
	}
	path := fmt.Sprintf("taggings/%d.json", taggingID)

	req, err := s.client.newRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DeleteTagging request: %w", err)
	}

	resp, err := s.client.do(req, nil) // Expects 204 No Content
	return resp, err
}
