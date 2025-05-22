package feedbinapi

import (
	"context"
	"fmt"
	"net/http"
)

// IconsService handles operations related to feed icons.
type IconsService struct {
	client *Client
}

// ListIcons retrieves feed icons (favicons) for all feeds the user is subscribed to.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/specs/content/icons.md#get-v2iconsjson
func (s *IconsService) ListIcons(ctx context.Context) ([]Icon, *Response, error) {
	path := "icons.json"
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ListIcons request: %w", err)
	}

	var icons []Icon
	resp, err := s.client.do(req, &icons)
	if err != nil {
		return nil, resp, err
	}
	return icons, resp, nil
}
