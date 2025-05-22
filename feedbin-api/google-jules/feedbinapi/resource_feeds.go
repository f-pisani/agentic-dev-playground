package feedbinapi

import (
	"context"
	"fmt"
	"net/http"
)

// FeedsService handles operations related to feeds.
type FeedsService struct {
	client *Client
}

// GetFeed retrieves a single feed by its ID.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/feeds.md#get-feed
func (s *FeedsService) GetFeed(ctx context.Context, feedID int64) (*Feed, *Response, error) {
	if feedID <= 0 {
		return nil, nil, fmt.Errorf("feedID must be a positive integer")
	}
	path := fmt.Sprintf("feeds/%d.json", feedID)

	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GetFeed request: %w", err)
	}

	var feed Feed
	resp, err := s.client.do(req, &feed)
	if err != nil {
		return nil, resp, err // Return response even on error for inspection
	}

	return &feed, resp, nil
}
