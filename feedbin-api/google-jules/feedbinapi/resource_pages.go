package feedbinapi

import (
	"context"
	"fmt"
	"net/http"
)

// PagesService handles operations related to creating entries from URLs (pages).
type PagesService struct {
	client *Client
}

// CreatePage creates a new entry from the URL of an article.
// The response, if successful, is the full Entry object.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/pages.md#post-v2pagesjson
func (s *PagesService) CreatePage(ctx context.Context, urlStr string, title *string) (*Entry, *Response, error) {
	if urlStr == "" {
		return nil, nil, fmt.Errorf("urlStr cannot be empty for CreatePage")
	}

	path := "pages.json"
	requestBody := &PageCreateRequest{URL: urlStr, Title: title}

	req, err := s.client.newRequest(ctx, http.MethodPost, path, nil, requestBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CreatePage request: %w", err)
	}

	var entry Entry // API returns a full Entry object on success
	resp, err := s.client.do(req, &entry)
	if err != nil {
		return nil, resp, err
	}
	return &entry, resp, nil
}
