package feedbin

import (
	"fmt"
	"net/http"
)

// FeedsService handles communication with the feed related
// methods of the Feedbin API.
type FeedsService service

// Get retrieves a single feed by its ID.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/feeds.md#get-feed
func (s *FeedsService) Get(feedID int64) (*Feed, *http.Response, error) {
	path := fmt.Sprintf("feeds/%d.json", feedID)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var feed Feed
	resp, err := s.client.do(req, &feed)
	if err != nil {
		return nil, resp, err
	}

	return &feed, resp, nil
}
