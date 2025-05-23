package feedbinapi

import (
	"net/http"
	// "github.com/google/go-querystring/query" // Not needed if no options
)

// IconsService handles operations related to feed icons.
type IconsService struct {
	client *Client
}

// NewIconsService creates a new service for icon operations.
func NewIconsService(client *Client) *IconsService {
	return &IconsService{client: client}
}

// IconListOptions specifies optional parameters for listing icons.
// The API doc for "GET /v2/icons.json" does not specify any query parameters.
type IconListOptions struct {
	// No options defined in spec for this endpoint.
}

// List retrieves all feed icons. Each icon object contains the host and base64 encoded image data.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/icons.md#get-icons
func (s *IconsService) List(opts *IconListOptions) ([]Icon, *http.Response, error) {
	path := "icons.json"
	// No options are currently supported by the endpoint, but including for future-proofing if any are added.
	// if opts != nil {
	// 	v, err := query.Values(opts)
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}
	// 	if params := v.Encode(); params != "" {
	// 		path = fmt.Sprintf("%s?%s", path, params)
	// 	}
	// }

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var icons []Icon
	resp, err := s.client.Do(req, &icons)
	if err != nil {
		return nil, resp, err
	}
	return icons, resp, nil
}
